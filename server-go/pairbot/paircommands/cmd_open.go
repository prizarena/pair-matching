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
)

const openCellCommandCode = "open"

var openCellCommand = bots.NewCallbackCommand(openCellCommandCode,
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		var player pairmodels.PairsPlayer
		q := callbackUrl.Query()
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

		for _, boardUserID := range board.UserIDs {
			if boardUserID != userID {
				playerEntityHolders = append(playerEntityHolders, &pairmodels.PairsPlayer{StringID: db.NewStrID(board.ID + ":" + boardUserID)})
			}
		}

		if len(playerEntityHolders) > 0 {
			if err = pairdal.DB.GetMulti(c, playerEntityHolders); err != nil {
				return
			}
		}

		player.ID = board.ID + ":" + userID

		err = pairdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			if err = pairdal.DB.Get(c, &player); err != nil && !db.IsNotFound(err) {
				return
			}
			if db.IsNotFound(err) {
				player.PairsPlayerEntity = &pairmodels.PairsPlayerEntity{
					Created: time.Now(),
				}
			}
			var players []pairmodels.PairsPlayer
			players = append(players, player)

			var changed bool
			if changed, err = pairgame.OpenCell(board.PairsBoardEntity, ca, player, players); err != nil {
				return
			} else if changed {
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
	if ca[0] < '0' || ca[0] > '9' {
		err = fmt.Errorf("unexpected Y value of '%v' parameter: %v", p, s)
		return
	}
	ca = turnbased.CellAddress(s)
	return
}
