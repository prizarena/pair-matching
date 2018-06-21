package pairgame

import (
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/prizarena/turn-based"
)

func OpenCell(board *pairmodels.PairsBoardEntity, ca turnbased.CellAddress, player pairmodels.PairsPlayer, players []pairmodels.PairsPlayer) (changed bool, err error) {
	if player.Open1 == "" {
		changed = true
		player.Open1 = ca
	} else if player.Open1 == ca && player.Open2 == "" {
		// changed = false
		return
	} else if player.Open2 == "" {
		changed = true
		alreadyOpened := board.GetCell(player.Open1)
		currentlyOpened := board.GetCell(ca)
		if alreadyOpened == currentlyOpened {
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
