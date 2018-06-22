package pairgame

import (
	"testing"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/strongo/db"
	"github.com/prizarena/turn-based"
)

func TestOpenCell(t *testing.T) {
	board := pairmodels.PairsBoardEntity{
		Size: "D3",
		Cells: "123456654321",
	}

	newPlayer := func(id string) pairmodels.PairsPlayer {
		return pairmodels.PairsPlayer{
			StringID: db.NewStrID(id),
			PairsPlayerEntity: &pairmodels.PairsPlayerEntity{

			},
		}
	}
	p1 := newPlayer("p1")
	p2 := newPlayer("p2")

	players := []pairmodels.PairsPlayer{p1, p2}
	playersByID := map[string]pairmodels.PairsPlayer {
		p1.ID: p1,
		p2.ID: p2,
	}

	type expectedPlayer struct {
		Open1, Open2 turnbased.CellAddress
		MatchedCount int
		MatchedItems string
	}

	type expected struct {
		changed bool
		players map[string]expectedPlayer
	}

	type Step struct {
		player pairmodels.PairsPlayer
		ca turnbased.CellAddress
		expected
	}
	var i int
	var step Step
	defer func(){
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic at step %v: %v", i+1, r)
		}
	}()
	for i, step = range []Step{
		{
			player: p1, ca: "A2",
			expected: expected{
				changed: true,
				players: map[string]expectedPlayer{
					"p1": {Open1: "A2"},
					"p2": {},
				},
			},
		},
		{
			player: p1, ca: "A3",
			expected: expected{
				changed: true,
				players: map[string]expectedPlayer{
					"p1": {Open1: "A2", Open2: "A3"},
					"p2": {},
				},
			},
		},
		{
			player: p1, ca: "B2",
			expected: expected{
				changed: true,
				players: map[string]expectedPlayer{
					"p1": {Open1: "B2", Open2: ""},
					"p2": {},
				},
			},
		},
		{
			player: p1, ca: "C2",
			expected: expected{
				changed: true,
				players: map[string]expectedPlayer{
					"p1": {Open1: "B2", Open2: "C2", MatchedItems: "6", MatchedCount: 1},
					"p2": {},
				},
			},
		},
		{
			player: p1, ca: "D3",
			expected: expected{
				changed: true,
				players: map[string]expectedPlayer{
					"p1": {Open1: "D3", Open2: "", MatchedItems: "6", MatchedCount: 1},
					"p2": {},
				},
			},
		},
	} {
		if changed, err := OpenCell(&board, step.ca, step.player, players); err != nil {
			t.Fatalf("Error at step #%v: %v", i, err)
		} else {
			if changed != step.expected.changed {
				if changed {
					t.Fatalf("Expected to be NOT changed at step #%v", i+1)
				} else {
					t.Fatalf("Expected to BE changed at step #%v", i+1)
				}
			}
			for pID, expectedP := range step.expected.players {
				actualP := playersByID[pID]
				var errorsCount int
				if actualP.Open1 != expectedP.Open1 {
					errorsCount++
					t.Errorf("Step #%v: actualP.Open1:%v != expectedP.Open1:%v", i+1, actualP.Open1, expectedP.Open1)
				}
				if actualP.Open2 != expectedP.Open2 {
					errorsCount++
					t.Errorf("Step #%v: actualP.Open2:%v != expectedP.Open2:%v", i+1, actualP.Open2, expectedP.Open2)
				}
				if actualP.MatchedItems != expectedP.MatchedItems {
					errorsCount++
					t.Errorf("Step #%v: actualP.MatchedItems:%v != expectedP.MatchedItems:%v", i+1, actualP.MatchedItems, expectedP.MatchedItems)
				}
				if actualP.MatchedCount != expectedP.MatchedCount {
					errorsCount++
					t.Errorf("Step #%v: actualP.MatchedCount:%v != expectedP.MatchedCount:%v", i+1, actualP.MatchedCount, expectedP.MatchedCount)
				}
				if errorsCount > 0 {
					t.Fatalf("Step #%v: errorsCount: %v", i+1, errorsCount)
				}
			}
		}
	}
}
