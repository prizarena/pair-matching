package pairmodels

import "github.com/strongo/db"

type BoardPlayerEntity struct {
	Turns   int    `datastore:",noindex,omitempty"`
	Matched string `datastore:",noindex,omitempty"`
	OpenCell    string `datastore:",noindex,omitempty"`
}

const BoardPlayerKind = "P"

type BoardPlayer struct {
	db.StringID
	*BoardPlayerEntity
}

var _ db.EntityHolder = (*BoardPlayer)(nil)

func (BoardPlayer) Kind() string {
	return BoardPlayerKind
}

func (player BoardPlayer) Entity() interface{} {
	return player.BoardPlayerEntity
}

func (BoardPlayer) NewEntity() interface{} {
	return new(BoardPlayerEntity)
}

func (player BoardPlayer) SetEntity(entity interface{}) {
	player.BoardPlayerEntity = entity.(*BoardPlayerEntity)
}