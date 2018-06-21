package pairmodels

import (
	"testing"
	"strings"
	"fmt"
)

func TestPairsBoardEntity_DrawBoard_ascii(t *testing.T) {
	board := PairsBoardEntity{
		Cells: "123456789abc",
		SizeX: 3,
		SizeY: 4,
	}
	expects := strings.Join([]string{"", "123", "456", "789", "abc", ""}, "\n")
	if result := board.DrawBoard(); result != expects {
		t.Error("Unexpected result:\n" + result)
	}
}

func TestPairsBoardEntity_DrawBoard_emoji(t *testing.T) {
	board := PairsBoardEntity{
		Cells: "🍇🍈🍉🍊🍋🍌🍍🍎🍏🍐🍑🍒",
		SizeX: 3,
		SizeY: 4,
	}
	expects := strings.Join([]string{"", "🍇🍈🍉", "🍊🍋🍌", "🍍🍎🍏", "🍐🍑🍒", ""}, "\n")
	if result := board.DrawBoard(); result != expects {
		t.Error("Unexpected result:\n" + result)
	}

	testShuffle := func(width, height int) {
		t.Helper()
		var board PairsBoardEntity
		board.SizeX = width
		board.SizeY = height
		board.Cells = Shuffle(width, height)
		rows := board.Rows()
		if len(rows) != height {
			t.Errorf("len(rows) != %v: %v", height, len(rows))
		}
		for rowIndex, row := range rows {
			if len(row) != width {
				t.Errorf("len(rows[%v]) != %v: %v", rowIndex, width, len(row))
			}
			for colIndex, r := range row {
				if r == 0 {
					t.Errorf("rows[%v][%v] == 0", colIndex, rowIndex)
				}
			}
		}
	}

	testShuffle(2,2)
	testShuffle(3,4)
	testShuffle(8,8)
}


func TestShuffle(t *testing.T) {

	test := func(n, x, y int) {
		s := Shuffle(x, y)
		if err := verifyBoard(x, y, s); err != nil {
			t.Errorf("Iteration %d shuffling %vx%v: %v", n, x, y, err)
		}
	}
	test(1,2, 2)
	test(2,3, 4)
	test(3,8, 8)
}

func verifyBoard(x, y int, s string) (err error){
	var itemsCount int
	counts := make(map[rune]int, x*y/2)
	for _, r := range s {
		itemsCount++
		counts[r]++
		if counts[r] > 2 {
			return fmt.Errorf("More then 2 items of %v", r)
		}

	}
	if itemsCount != x*y {
		fmt.Errorf("Expectet %v items, got %v", x*y, itemsCount)
	}
	return nil
}

func TestGetCell(t *testing.T) {
	board := PairsBoardEntity{
		Cells: "🍇🍈🍉🍊🍋🍌🍍🍎🍏🍐🍑🍒",
		SizeX: 3,
		SizeY: 4,
	}
	testCell := func(x, y int, expects rune) {
		t.Helper()
		if v := board.GetCell(x, y); v != expects {
			t.Errorf("%d:%d expected %v got %v", x, y, string(expects), string(v))
		}
	}
	testCell(1, 1, '🍇')
	testCell(2, 1, '🍈')
	testCell(3, 1, '🍉')

	testCell(1, 2, '🍊')
	testCell(2, 2, '🍋')
	testCell(3, 2, '🍌')

	testCell(1, 3, '🍍')
	testCell(2, 3, '🍎')
	testCell(3, 3, '🍏')

	testCell(1, 4, '🍐')
	testCell(2, 4, '🍑')
	testCell(3, 4, '🍒')
}