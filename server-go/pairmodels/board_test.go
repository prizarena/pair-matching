package pairmodels

import (
	"testing"
	"strings"
)

func TestPairsBoardEntity_DrawBoard(t *testing.T) {
	board := PairsBoardEntity{
		Cells: "123456789abc",
		SizeX: 3,
		SizeY: 4,
	}
	expects := strings.Join([]string{"123", "456", "789", "abc"}, "\n") + "\n"
	if result := board.DrawBoard(); result != expects {
		t.Error("Unexpected result:\n" + result)
	}
}
