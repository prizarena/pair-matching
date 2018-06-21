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
	"github.com/strongo/slices"
)

const newBoardCommandCode = "new"

func newBoardCallbackData(width, height int) string {
	return fmt.Sprintf("new?w=%d&h=%v", width, height)
}

var newBoardCommand = bots.NewCallbackCommand(
	newBoardCommandCode,
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		q := callbackUrl.Query()
		var widht, height int
		if widht, height, err = getPoint(q, "w", "h"); err != nil {
			return
		}

		var board pairmodels.PairsBoard
		board.ID = whc.Input().(telegram.TgWebhookCallbackQuery).GetInlineMessageID()

		userID := whc.AppUserStrID()
		err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			if err = pairdal.DB.Get(tc, &board); err != nil && !db.IsNotFound(err) {
				return
			}
			var changed bool
			if db.IsNotFound(err) {
				changed = true
				board = pairmodels.PairsBoard{
					PairsBoardEntity: &pairmodels.PairsBoardEntity{
						SizeY: widht,
						SizeX: height,
						Cells: pairmodels.Shuffle(widht, height),
						BoardEntityBase: turnbased.BoardEntityBase{
							UserIDs: []string{},
						},
					},
				}
			} else if slices.IsInStringSlice(userID, board.UserIDs) {
				changed = true
				board.UserIDs = append(board.UserIDs, userID)
			}
			if changed {
				if err = pairdal.DB.Update(tc, &board); err != nil {
					return
				}
			}
			return
		}, db.SingleGroupTransaction)
		if err != nil {
			return
		}
		// TODO: check and notify if another user already selected different board size.
		m, err = renderPairsBoardMessage(whc, nil, board, nil)
		return
	},
)
