package pairgame

import (
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/prizarena/turn-based"
	"github.com/pkg/errors"
	"bytes"
	"strings"
	)

var (
	ErrAlreadyMatched   = errors.New("already matched")
	ErrBoardIsCompleted = errors.New("board is already completed")
)

func OpenCell(
	board *pairmodels.PairsBoardEntity, ca turnbased.CellAddress, player pairmodels.PairsPlayer, players []pairmodels.PairsPlayer,
) (
	changed bool, changedPlayers []pairmodels.PairsPlayer, err error,
) {
	if ca == "" {
		panic("Cell address is required to open a cell")
	}
	if board.IsCompleted(players) {
		err = ErrBoardIsCompleted
		return
	}
	currentlyOpened := board.GetCell(ca)

	if player.Open1 != "" && player.Open2 != "" {
		// If player has 2 tiles opened close them before opening next one
		changed = true
		player.Open1 = ""
		player.Open2 = ""
	}

	playerChanged := func(p pairmodels.PairsPlayer) {
		for _, cp := range changedPlayers {
			if cp.ID == p.ID {
				return
			}
		}
		changedPlayers = append(changedPlayers, p)
	}

	{ // Close already matched cells
		allAlreadyMatchedItems := func() string {
			var s bytes.Buffer
			for _, p := range players {
				s.WriteString(p.MatchedItems)
			}
			return s.String()
		}()
		closeAlreadyMatched := func(p pairmodels.PairsPlayer) (pChanged bool) {
			if p.Open1 != "" && strings.Contains(allAlreadyMatchedItems, string(board.GetCell(p.Open1))) {
				p.Open1 = ""
				pChanged = true
			}
			if p.Open2 != "" && strings.Contains(allAlreadyMatchedItems, string(board.GetCell(p.Open2))) {
				p.Open2 = ""
				pChanged = true
			}
			return
		}
		changed = closeAlreadyMatched(player) || changed
		for _, p := range players {
			if closeAlreadyMatched(p) {
				changedPlayers = append(changedPlayers, p)
			}
		}
	}

	atLeastOneMatched := false

	for _, p := range players {
		if strings.Contains(p.MatchedItems, string(currentlyOpened)) {
			err = ErrAlreadyMatched
			if p.Open1 != "" && board.GetCell(p.Open1) == currentlyOpened {
				p.Open1 = ""
				playerChanged(p)
			}
			if p.Open2 != "" && board.GetCell(p.Open2) == currentlyOpened {
				p.Open2 = ""
				playerChanged(p)
			}
			return
		}

		match := func(openN int, openField turnbased.CellAddress) (matched bool) {
			if matched = board.GetCell(openField) == currentlyOpened; matched {
				playerChanged(p)
			 	// if p.ID == player.ID {
				// 	// TODO: Unit tests VERY needed!
				// 	// Current player opened a 3d tile that is 1 of 2 previously opened by him/here.
				// 	// This is not a real match.
				// 	matched = false
				// } else {
				// 	err = fmt.Errorf("unexpected player state - both Open1 and Open2 are not empty: playerID=%v, Open1=%v, Open2=%v, ca=%v",
				// 		player.ID, player.Open1, player.Open2, ca)
				// 	return
				// }
				// switch openN {
				// case 1:
				// 	p.Open1 = p.Open2
				// 	p.Open2 = ""
				// case 2:
				// 	p.Open2 = ""
				// default:
				// 	panic("openN != 1|2")
				// }
			}
			return
		}

		isMatched := p.Open1 != "" && match(1, p.Open1)
		if err != nil {
			return
		}
		if !isMatched {
			isMatched = p.Open2 != "" && match(2, p.Open2)
			if err != nil {
				return
			}
		}

		if isMatched {
			atLeastOneMatched = true
			// Do not break as theoretically in case of race could match more than 1.
		}
	}

	if atLeastOneMatched {
		changed = true
		if player.Open1 == "" {
			player.Open1 = ca
		} else if player.Open2 == "" {
			player.Open2 = ca
		}
		player.MatchedCount++
		player.MatchedItems += string(currentlyOpened)
		return
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
