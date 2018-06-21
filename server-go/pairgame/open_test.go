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
				players: map[string]expectedPlayer{
					"p1": {Open1: "A2"},
					"p2": {},
				},
			},
		},
		{
			player: p1, ca: "A3",
			expected: expected{
				players: map[string]expectedPlayer{
					"p1": {Open1: "A2", Open2: "A3"},
					"p2": {},
				},
			},
		},
		{
			player: p1, ca: "B2",
			expected: expected{
				players: map[string]expectedPlayer{
					"p1": {Open1: "B2", Open2: ""},
					"p2": {},
				},
			},
		},
	} {
		if changed, err := OpenCell(&board, step.ca, step.player, players); err != nil {
			t.Fatalf("step #%v: %v", i, err)
			if changed != step.expected.changed {
				t.Fatalf("step #%v", i)
			}
			for pID, expectedP := range step.expected.players {
				actualP := playersByID[pID]
				if actualP.Open1 != expectedP.Open1 {
					t.Errorf("actualP.Open1 != expectedP.open1")
				}
				if actualP.Open2 != expectedP.Open2 {
					t.Errorf("actualP.Open2 != expectedP.open2")
				}
				if actualP.MatchedItems != expectedP.MatchedItems {
					t.Errorf("actualP.MatchedItems != expectedP.MatchedItems")
				}
			}
		}
	}
}
