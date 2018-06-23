package pairgame

import (
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/prizarena/turn-based"
	"github.com/pkg/errors"
	"strings"
	"bytes"
	"fmt"
)

var ErrAlreadyMatched = errors.New("already matched")

func OpenCell(board *pairmodels.PairsBoardEntity, ca turnbased.CellAddress, player pairmodels.PairsPlayer, players []pairmodels.PairsPlayer) (changed bool, err error) {
	if ca == "" {
		panic("Cell address is required to open a cell")
	}
	currentlyOpened := board.GetCell(ca)

	{ // Close already matched cells
		allAlreadyMatchedItems := func() string {
			var s bytes.Buffer
			for _, p := range players {
				s.WriteString(p.MatchedItems)
			}
			return s.String()
		}()
		closeAlreadyMatched := func(p pairmodels.PairsPlayer) {
			if p.Open1 != "" && strings.Contains(allAlreadyMatchedItems, string(board.GetCell(p.Open1))) {
				p.Open1 = ""
			}
			if p.Open2 != "" && strings.Contains(allAlreadyMatchedItems, string(board.GetCell(p.Open2))) {
				p.Open2 = ""
			}
		}
		closeAlreadyMatched(player)
		for _, p := range players {
			closeAlreadyMatched(p)
		}
	}

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

		match := func(openField *turnbased.CellAddress) (matched bool) {
			if matched = board.GetCell(*openField) == currentlyOpened; matched {
				player.MatchedCount++
				player.MatchedItems += string(currentlyOpened)
				if player.Open1 == "" {
					player.Open1 = ca
				} else if player.Open2 == "" {
					player.Open2 = ca
				} else {
					err = fmt.Errorf("unexpected player state - both Open1 and Open2 are not empty: playerID=%v, Open1=%v, Open2=%v, ca=%v",
						player.ID, player.Open1, player.Open2, ca)
					return
				}
				changed = true
			}
			return
		}

		isMatched := p.Open1 != "" && match(&p.Open1)
		if err != nil {
			return
		}
		if !isMatched {
			isMatched = p.Open2 != "" && match(&p.Open2)
			if err != nil {
				return
			}
		}
		if isMatched {
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
		player.Open2 = ca

		// This check is just in case. Actually should be caught by match() function above.
		// if alreadyOpened := board.GetCell(player.Open1); alreadyOpened == currentlyOpened {
		// 	player.Open2 = ca
		// 	player.MatchedCount++
		// 	player.MatchedItems += string(currentlyOpened)
		// 	return
		// }
	} else if player.Open1 != "" && player.Open2 != "" {
		changed = true
		player.Open1 = ca
		player.Open2 = ""
	} else {
		panic("should not be here")
	}
	return
}
