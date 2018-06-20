package pairmodels

import (
	"github.com/strongo/db"
	"bytes"
)

type PairsBoardEntity struct {
	Cells string `datastore:",noindex,omitempty"`
	SizeX int
	SizeY int
}

const PairBoardKind = "B"

type PairsBoard struct {
	db.StringID
	*PairsBoardEntity
}

var _ db.EntityHolder = (*BoardPlayer)(nil)

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

func (board PairsBoardEntity) DrawBoard() string {
	s := new(bytes.Buffer)

	s.WriteRune('\n')
	for i := 0; i < board.SizeY; i++ {
		first := board.SizeX*i
		s.WriteString(board.Cells[first:first+board.SizeX])
		s.WriteRune('\n')
	}
	return s.String()
}
