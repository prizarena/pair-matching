package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"strconv"
	"github.com/pkg/errors"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/prizarena/pair-matching/server-go/pairdal"
	"context"
	"github.com/strongo/db"
	"time"
)

const openCellCommandCode = "open"

var openCellCommand = bots.NewCallbackCommand(openCellCommandCode,
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		var player pairmodels.PairsPlayer
		var x, y int
		q := callbackUrl.Query()
		if x, y, err = getPoint(q, "x", "y"); err != nil {
			return
		}

		var board pairmodels.PairsBoard
		board.ID = q.Get("b")
		if err = pairdal.DB.Get(c, &board); err != nil {
			return
		}

		userID := whc.AppUserStrID()

		player.ID = board.ID + ":" + userID

		err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			if err = pairdal.DB.Get(c, &player); err != nil && !db.IsNotFound(err) {
				return
			}
			var changed bool
			if db.IsNotFound(err) {
				player.PairsPlayerEntity = &pairmodels.PairsPlayerEntity{
					Created: time.Now(),
				}
			}
			if player.OpenX == 0 && player.OpenY == 0 {
				changed = true
				player.OpenX = x
				player.OpenY = y
			} else {
				changed = true
				alreadyOpened := board.GetCell(player.OpenX, player.OpenY)
				currentlyOpened := board.GetCell(x, y)
				if alreadyOpened == currentlyOpened {
					player.MatchedCount++
					player.MatchedItems += string(currentlyOpened)
				}
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

		// if player.OpenX == x && player.OpenY == y {
		// 	m, err = renderPairsBoardMessage(whc, nil, board, []pairmodels.PairsPlayer{player})
		// 	return
		// }

		m, err = renderPairsBoardMessage(whc, nil, board, []pairmodels.PairsPlayer{player})
		return
	},
)

func getPoint(v url.Values, p1, p2 string) (v1, v2 int, err error) {
	if v1, err = strconv.Atoi(v.Get(p1)); err != nil {
		err = errors.WithMessage(err, "invalid "+p1)
		return
	}
	if v2, err = strconv.Atoi(v.Get(p2)); err != nil {
		err = errors.WithMessage(err, "invalid "+p2)
		return
	}
	return
}
