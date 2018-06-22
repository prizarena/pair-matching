package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/prizarena/pair-matching/server-go/pairdal"
	"context"
	"github.com/strongo/db"
	"time"
	"github.com/prizarena/pair-matching/server-go/pairgame"
	"github.com/prizarena/turn-based"
	"fmt"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/log"
	"github.com/strongo/slices"
)

const openCellCommandCode = "open"

func openCellCallbackData(ca turnbased.CellAddress, boardID, lang string) string {
	return fmt.Sprintf(openCellCommandCode+"?c=%v&b=%v&l=%v", ca, boardID, lang)
}

var openCellCommand = bots.NewCallbackCommand(openCellCommandCode,
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()

		q := callbackUrl.Query()
		if err = whc.SetLocale(q.Get("l")); err != nil {
			return
		}
		var player pairmodels.PairsPlayer
		var ca turnbased.CellAddress
		if ca, err = getCellAddress(q, "c"); err != nil {
			return
		}

		var board pairmodels.PairsBoard
		board.ID = q.Get("b")
		if err = pairdal.DB.Get(c, &board); err != nil {
			return
		}

		userID := whc.AppUserStrID()

		playerEntityHolders := make([]db.EntityHolder, 0, len(board.UserIDs)+1)

		isNewUser := true
		for _, boardUserID := range board.UserIDs {
			if boardUserID == userID {
				isNewUser = false
			} else {
				playerEntityHolders = append(playerEntityHolders, &pairmodels.PairsPlayer{StringID: db.NewStrID(board.ID + ":" + boardUserID)})
			}
		}

		var userName string
		if isNewUser {
			var botAppUser bots.BotAppUser
			if botAppUser, err = whc.GetAppUser(); err != nil {
				return
			}
			err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
				if err = pairdal.DB.Get(c, &board); err != nil {
					return
				}
				if !slices.IsInStringSlice(userID, board.UserIDs) {
					userName = botAppUser.(*pairmodels.UserEntity).FullName()
					board.AddUser(userID, userName)
					if err = pairdal.DB.Update(c, &board); err != nil {
						return
					}
				}
				return
			}, db.CrossGroupTransaction)
			if err != nil {
				return
			}
		} else {
			for i, uID := range board.UserIDs {
				if uID == userID {
					userName = board.UserNames[i]
				}
			}
		}

		if len(playerEntityHolders) > 0 {
			if err = pairdal.DB.GetMulti(c, playerEntityHolders); err != nil {
				return
			}
		}

		player.ID = board.ID + ":" + userID

		players := make([]pairmodels.PairsPlayer, len(board.UserIDs))

		var isAlreadyMatched bool
		// =[ Start of transaction ]=
		err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			if err = pairdal.DB.Get(c, &player); err != nil && !db.IsNotFound(err) {
				return
			}
			if db.IsNotFound(err) {
				player.PairsPlayerEntity = &pairmodels.PairsPlayerEntity{
					Created: time.Now(),
					UserName: userName,
				}
			}

			// Populate players in same order as joined board to keep displaying in same order.
			for i, uID := range board.UserIDs {
				if player.UserID() == uID {
					players[i] = player
				} else {
					for _, eh := range playerEntityHolders {
						if p := *eh.(*pairmodels.PairsPlayer); p.UserID() == uID {
							players[i] = p
						}
					}
				}
			}

			var changed bool
			// ===================================================================================================
			if changed, err = pairgame.OpenCell(board.PairsBoardEntity, ca, player, players); err != nil {
				// ================================================================================================
				if err == pairgame.ErrAlreadyMatched {
					isAlreadyMatched = true
					err = nil
				} else {
					return
				}
			}
			if userName != "" && player.UserName != userName {
				player.UserName = userName
				changed = true
			}
			if changed {
				if err = pairdal.DB.Update(c, &player); err != nil && !db.IsNotFound(err) {
					return
				}
			}
			return
		}, db.SingleGroupTransaction)
		if err != nil {
			return
		}
		// =[ End of transaction ]=

		// if player.OpenX == x && player.OpenY == y {
		// 	m, err = renderPairsBoardMessage(whc, nil, board, []pairmodels.PairsPlayer{player})
		// 	return
		// }
		if isAlreadyMatched {
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
				Text:      "This cell is already matched",
				CacheTime: 10,
			})
			if _, err = whc.Responder().SendMessage(c, m, bots.BotAPISendMessageOverHTTPS); err != nil {
				log.Errorf(c, "Failed to send already matched alert: %v", err)
				err = nil
			}
		}
		return renderPairsBoardMessage(whc, nil, board, players)
	},
)

func getSize(v url.Values, p string) (size turnbased.Size, err error) {
	var ca turnbased.CellAddress
	if ca, err = getCellAddress(v, p); err != nil {
		return
	}
	return turnbased.Size(ca), nil
}
func getCellAddress(v url.Values, p string) (ca turnbased.CellAddress, err error) {
	s := v.Get(p)
	if len(s) != 2 {
		err = fmt.Errorf("unexpected length of '%v' parameter: %v", p, len(s))
		return
	}
	if s[0] < 'A' || s[0] > 'Z' {
		err = fmt.Errorf("unexpected X value of '%v' parameter: %v", p, s)
		return
	}
	if s[1] < '0' || s[1] > '9' {
		err = fmt.Errorf("unexpected Y value of '%v' parameter: %v", p, s)
		return
	}
	ca = turnbased.CellAddress(s)
	return
}
