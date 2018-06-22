package pairgame

import (
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/prizarena/turn-based"
		"github.com/pkg/errors"
	"strings"
)

var ErrAlreadyMatched = errors.New("already matched")

func OpenCell(board *pairmodels.PairsBoardEntity, ca turnbased.CellAddress, player pairmodels.PairsPlayer, players []pairmodels.PairsPlayer) (changed bool, err error) {
	if ca == "" {
		panic("Cell address is required to open a cell")
	}
	currentlyOpened := board.GetCell(ca)
	for _, p := range players {
		if strings.Contains(p.MatchedItems, string(currentlyOpened)) {
			err = ErrAlreadyMatched
			if p.Open1 != "" && board.GetCell(p.Open1) == currentlyOpened {
				changed = true
				p.Open1 = ""
			}
			if p.Open2 != "" && board.GetCell(p.Open2) == currentlyOpened {
				changed = true
				p.Open2 = ""
			}
			return
		}
	}

	if player.Open1 == "" {
		changed = true
		player.Open1 = ca
	} else if player.Open1 == ca && player.Open2 == "" {
		// changed = false
		return
	} else if player.Open2 == "" {

		changed = true
		alreadyOpened := board.GetCell(player.Open1)

		if alreadyOpened == currentlyOpened {
			player.Open2 = ca
			player.MatchedCount++
			player.MatchedItems += string(currentlyOpened)
			return
		}
		player.Open2 = ca
	} else if player.Open1 != "" && player.Open2 != "" {
		changed = true
		player.Open1 = ca
		player.Open2 = ""
	} else {
		panic("should not be here")
	}
	return
}
