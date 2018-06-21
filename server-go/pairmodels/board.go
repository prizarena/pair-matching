package pairmodels

import (
	"github.com/strongo/db"
	"bytes"
)

type PairsBoardEntity struct {
	Cells     string `datastore:",noindex,omitempty"`
	SizeX     int
	SizeY     int
	UserIDs   []string
	UserNames []string
}

const PairBoardKind = "B"

type PairsBoard struct {
	db.StringID
	*PairsBoardEntity
}

var _ db.EntityHolder = (*PairsPlayer)(nil)

func (PairsBoard) Kind() string {
	return PairBoardKind
}

func (eh PairsBoard) Entity() interface{} {
	return eh.PairsBoardEntity
}

func (PairsBoard) NewEntity() interface{} {
	return new(PairsBoardEntity)
}

func (eh PairsBoard) SetEntity(entity interface{}) {
	eh.PairsBoardEntity = entity.(*PairsBoardEntity)
}

func (board PairsBoardEntity) Rows() (rows [][]rune) {
	var x, y = 0, 0
	rows = make([][]rune, board.SizeY)
	if board.SizeX == 0 {
		return
	}
	rows[0] = make([]rune, board.SizeX)
	for _, r := range board.Cells {
		rows[y][x] = r
		if x++; x == board.SizeX {
			x = 0
			if y++; y < board.SizeY {
				rows[y] = make([]rune, board.SizeX)
			}
		}
	}
	return
}

func (board PairsBoardEntity) DrawBoard() string {
	s := new(bytes.Buffer)

	s.WriteRune('\n')
	rows := board.Rows()
	for _, row := range rows {
		for _, r := range row {
			s.WriteRune(r)
		}
		s.WriteRune('\n')
	}
	return s.String()
}
