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
	"strconv"
	"github.com/pkg/errors"
	"bytes"
	"github.com/prizarena/prizarena-public/pamodels"
)

const openCellCommandCode = "open"

func openCellCallbackData(ca turnbased.CellAddress, playersCount int, boardID, userID, lang string) string {
	s := new(bytes.Buffer)
	s.WriteString(openCellCommandCode)
	fmt.Fprintf(s, "?c=%v&p=%v&l=%v", ca, playersCount, lang)
	// fmt.Fprintf(s, "?c=%v&p=%v&b=%v&l=%v", ca, playersCount, boardID, lang)
	if playersCount == 1 {
		s.WriteString("&u=" + userID)
	}
	if boardID != ""  {
		s.WriteString("&b="+boardID)
	}
	return s.String()
}

var openCellCommand = bots.NewCallbackCommand(openCellCommandCode,
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		log.Debugf(c, "openCellCommand.CallbackAction()")

		q := callbackUrl.Query()
		if err = whc.SetLocale(q.Get("l")); err != nil {
			return
		}
		var (
			player              pairmodels.PairsPlayer
			ca                  turnbased.CellAddress
			playersCount        int
			board               pairmodels.PairsBoard
			playerEntityHolders []db.EntityHolder
		)

		if ca, err = getCellAddress(q, "c"); err != nil {
			err = errors.WithMessage(err, "parameter 'c' is not a valid cell address")
			return
		}

		if playersCount, err = strconv.Atoi(q.Get("p")); err != nil {
			err = errors.WithMessage(err, "parameter 'p' (e.g. playersCount) is not an integer")
			return
		} else if playersCount < 0 {
			err = fmt.Errorf("bad request: playersCount expected to be 0 or greater, got: %v", playersCount)
			return
		}

		if board.ID = q.Get("b"); board.ID == "" {
			if tgCallbackQuery, ok := whc.Input().(telegram.TgWebhookCallbackQuery); ok {
				board.ID = tgCallbackQuery.GetInlineMessageID()
			}
			if board.ID == "" {
				err = errors.New("bad request: parameter 'b' e.g. board ID is required")
				return
			}
		}

		userID := whc.AppUserStrID()

		isNewUser := playersCount == 0 || q.Get("u") != userID

		if playersCount > 1 {
			if err = pairdal.DB.Get(c, &board); err != nil {
				return
			}
			playerEntityHolders = make([]db.EntityHolder, 0, len(board.UserIDs)+1)
			for _, boardUserID := range board.UserIDs {
				if boardUserID == userID {
					isNewUser = false
				} else {
					playerEntityHolders = append(playerEntityHolders, &pairmodels.PairsPlayer{StringID: db.NewStrID(board.ID + ":" + boardUserID)})
				}
			}
		} else if isNewUser {
			playerEntityHolders = make([]db.EntityHolder, 0, playersCount+1)
		} else {
			playerEntityHolders = make([]db.EntityHolder, 0, 1)
		}

		var userName string

		var addUserToBoardCalled int
		addUserToBoard := func() (err error) {
			if addUserToBoardCalled++; addUserToBoardCalled > 1 {
				err = errors.New("addUserToBoardCalled should be called just once")
				return
			}
			log.Debugf(c, "addUserToBoard")
			var botAppUser bots.BotAppUser
			if !slices.IsInStringSlice(userID, board.UserIDs) {
				if userName == "" {
					if botAppUser, err = whc.GetAppUser(); err != nil {
						return
					}
					userName = botAppUser.(*pairmodels.UserEntity).FullName()
				}
				board.AddUser(userID, userName)
			}
			return
		}

		if playersCount > 1 {
			if isNewUser {
				err = pairdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
					if err = pairdal.DB.Get(c, &board); err != nil {
						return
					}
					if err = addUserToBoard(); err != nil {
						return
					}
					if err = pairdal.DB.Update(c, &board); err != nil {
						return
					}
					return
				}, db.SingleGroupTransaction)
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
		}

		if len(playerEntityHolders) > 0 {
			if err = pairdal.DB.GetMulti(c, playerEntityHolders); err != nil {
				return
			}
		}

		player.ID = board.ID + ":" + userID

		var players []pairmodels.PairsPlayer

		var isAlreadyMatched, isAlreadyOpen bool
		// =[ Start of transaction ]=
		txOptions := db.CrossGroupTransaction
		if playersCount > 1 {
			txOptions = db.SingleGroupTransaction
		} else {
			txOptions = db.CrossGroupTransaction
		}

		err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {

			newPlayerEntity := func() *pairmodels.PairsPlayerEntity {
				return &pairmodels.PairsPlayerEntity{
					PlayerCreated: time.Now(),
					UserName:      userName,
				}
			}

			var firstPlayer pairmodels.PairsPlayer

			if playersCount < 2 {
				if err = pairdal.DB.Get(c, &board); err != nil && !db.IsNotFound(err) {
					return
				}
				isNewUser = true
				for i, uID := range board.UserIDs {
					if uID == userID {
						isNewUser = false
						if userName == "" {
							userName = board.UserNames[i]
						}
						break
					}
				}
				log.Debugf(c, "isNewUser=%v, len(board.UserIDs)=%v", isNewUser, len(board.UserIDs))
				if isNewUser {
					if len(board.UserIDs) == 1 {
						firstPlayer.ID = pairmodels.NewPlayerID(board.ID, board.UserIDs[0])
						if err = pairdal.DB.Get(c, &firstPlayer); err == nil { // Let's lock the first player entity
							log.Warningf(c, "Unexpected but not critical: first user has PairsPlayer entity")
						} else if db.IsNotFound(err) { // This is the expected case
							err = nil
							entityCopy := board.PairsPlayerEntity
							entityCopy.UserName = board.UserNames[0]
							firstPlayer.PairsPlayerEntity = &entityCopy
						} else { // some real error returned
							return
						}
						board.PairsPlayerEntity = pairmodels.PairsPlayerEntity{} // Cleanup board from single player
						playerEntityHolders = append(playerEntityHolders, &firstPlayer)
					}
					addUserToBoard()
				}

				switch len(board.UserIDs) {
				case 0:
					player.PairsPlayerEntity = newPlayerEntity()
				case 1:
					if board.UserIDs[0] == userID {
						player.PairsPlayerEntity = &board.PairsPlayerEntity
						player.UserName = board.UserNames[0]
					}
				}
			}

			if player.PairsPlayerEntity == nil {
				if err = pairdal.DB.Get(c, &player); err != nil && !db.IsNotFound(err) {
					return
				}
				if db.IsNotFound(err) {
					err = nil
					log.Debugf(c, "we got a new user for the board, so create player.PairsPlayerEntity")
					player.ID = pairmodels.NewPlayerID(board.ID, userID)
					player.PairsPlayerEntity = newPlayerEntity()
				}
			}

			// Populate players in same order as joined board to keep displaying in same order.
			{
				players = make([]pairmodels.PairsPlayer, len(board.UserIDs))
				var playersSetCount int
				for i, uID := range board.UserIDs {
					if uID == userID {
						players[i] = player
						log.Debugf(c, "uID == userID: %v", uID)
						playersSetCount++
					} else {
						for _, eh := range playerEntityHolders {
							if p := eh.(*pairmodels.PairsPlayer); p.UserID() == uID {
								log.Debugf(c, "p.UserID() == uID: %v", uID)
								players[i] = *p
								playersSetCount++
							}
						}
					}
				}
				if playersSetCount != len(players) {
					err = fmt.Errorf("playersSetCount != len(players): %v != %v, len(playerEntityHolders)=%v, player.ID=%v, board.UserIDs=%v", playersSetCount, len(players), len(playerEntityHolders), player.ID, board.UserIDs)
					return
				}
			}

			var changed bool
			// var changedPlayers []pairmodels.PairsPlayer
			// ===================================================================================================
			if changed, _, err = pairgame.OpenCell(board.PairsBoardEntity, ca, player, players); err != nil {
				// ================================================================================================
				switch err {
				case pairgame.ErrAlreadyOpen:
					isAlreadyOpen = true
					err = nil
				case pairgame.ErrAlreadyMatched:
					isAlreadyMatched = true
					err = nil
				case pairgame.ErrBoardIsCompleted:
					log.Debugf(c, err.Error())
					err = nil
				}
			}
			if userName != "" && player.UserName != userName {
				player.UserName = userName
				changed = true
			}
			if changed {
				if playersCount < 2 && len(players) == 1 {
					board.PairsPlayerEntity.UserName = ""
					if err = pairdal.DB.Update(c, &board); err != nil && !db.IsNotFound(err) {
						return
					}
					board.PairsPlayerEntity.UserName = board.UserNames[0]
				} else if playersCount < 2 && len(players) > 1 {
					var entitiesToUpdate []db.EntityHolder
					if firstPlayer.ID == "" {
						entitiesToUpdate = []db.EntityHolder{&board, &player}
					} else {
						entitiesToUpdate = []db.EntityHolder{&board, &player, &firstPlayer}
					}
					if err = pairdal.DB.UpdateMulti(c, entitiesToUpdate); err != nil && !db.IsNotFound(err) {
						return
					}
				} else if playersCount > 1 && len(players) > 1 {
					if err = pairdal.DB.Update(c, &player); err != nil && !db.IsNotFound(err) {
						return
					}
				} else {
					err = fmt.Errorf("reached unexpected branch: playersCount=%v && len(players)=%v", playersCount, len(players))
					return
				}
			}
			return
		}, txOptions)
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
		if isAlreadyOpen {
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
				Text:      "This cell is already open by you",
				CacheTime: 10,
			})
			if _, err = whc.Responder().SendMessage(c, m, bots.BotAPISendMessageOverHTTPS); err != nil {
				log.Errorf(c, "Failed to send already matched alert: %v", err)
				err = nil
			}
		}
		tournament := pamodels.Tournament{StringID: db.NewStrID(board.TournamentID)}
		return renderPairsBoardMessage(whc, tournament, board, userID, players)
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
