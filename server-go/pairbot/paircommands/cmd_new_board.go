package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"fmt"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/prizarena/pair-matching/server-go/pairdal"
	"context"
	"github.com/strongo/db"
	"github.com/prizarena/turn-based"
		"github.com/strongo/log"
	"time"
	"bytes"
	"strconv"
)

const newBoardCommandCode = "new"

func getNewBoardCallbackData(width, height, maxUsersLimit int, tournamentID, lang string) string {
	s := new(bytes.Buffer)
	fmt.Fprintf(s,"new?s=%v&l=%v", turnbased.NewSize(width, height), lang)
	if tournamentID != "" {
		fmt.Fprint(s, "&t=" + tournamentID)
	}
	if maxUsersLimit > 0 {
		fmt.Fprint(s, "&max=" + strconv.Itoa(maxUsersLimit))
	} else if maxUsersLimit < 0 {
		panic(fmt.Sprintf("maxUsersLimit < 0: %v", maxUsersLimit))
	}
	return s.String()
}

var newBoardCommand = bots.NewCallbackCommand(
	newBoardCommandCode,
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		q := callbackUrl.Query()
		var size turnbased.Size
		if size, err = getSize(q, "s"); err != nil {
			return
		}

		var maxUsersLimit int
		if s := q.Get("max"); s != "" {
			if maxUsersLimit, err = strconv.Atoi(s); err != nil {
				return
			}
		}

		if err = whc.SetLocale(q.Get("l")); err != nil {
			return
		}

		var board pairmodels.PairsBoard
		board.ID = whc.Input().(telegram.TgWebhookCallbackQuery).GetInlineMessageID()

		userID := whc.AppUserStrID()
		// var botAppUser bots.BotAppUser
		// if botAppUser, err = whc.GetAppUser(); err != nil {
		// 	return
		// }
		err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			if err = pairdal.DB.Get(tc, &board); err != nil && !db.IsNotFound(err) {
				return
			}
			var changed bool
			if err == nil { // Existing entity
				log.Debugf(c, "Existing board entity")
				if boardUsersCount := len(board.UserIDs); boardUsersCount > 1 {
					log.Debugf(c, "Will delete %v player entities", boardUsersCount)
					players := make([]db.EntityHolder, boardUsersCount)
					for i, userID := range board.UserIDs {
						players[i] = &pairmodels.PairsPlayer{StringID: db.NewStrID(pairmodels.NewPlayerID(board.ID, userID))}
					}
					if err = pairdal.DB.DeleteMulti(tc, players); err != nil {
						return
					}
				}
				now := time.Now()
				if board.Created.Before(now.Add(-time.Second*2)) {
					board.Created = now
					board.Size = size
					board.Cells = pairmodels.NewCells(size.Width(), size.Height())
					board.PairsPlayerEntity = pairmodels.PairsPlayerEntity{}
					board.UserIDs = []string{}
					board.UserNames = []string{}
					board.UserWins = []int{}
					board.UsersMax = maxUsersLimit
					changed = true
				}
			} else if db.IsNotFound(err) {
				log.Debugf(c, "New board entity")
				changed = true
				board.PairsBoardEntity = &pairmodels.PairsBoardEntity{
					BoardEntityBase: turnbased.BoardEntityBase{
						Created: time.Now(),
						CreatorUserID: userID,
					},
					Size:  size,
					Cells: pairmodels.NewCells(size.Width(), size.Height()),
				}
			}
			// if !slices.IsInStringSlice(userID, board.UserIDs) {
			// 	changed = true
			// 	board.AddUser(userID, botAppUser.(*pairmodels.UserEntity).FullName())
			// }
			if changed {
				if err = pairdal.DB.Update(tc, &board); err != nil {
					return
				}
			}
			return
		}, db.CrossGroupTransaction)
		if err != nil {
			return
		}
		// TODO: check and notify if another user already selected different board size.
		m, err = renderPairsBoardMessage(whc, nil, board, userID,nil)
		return
	},
)
