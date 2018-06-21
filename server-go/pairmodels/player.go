package pairmodels

import (
	"github.com/strongo/db"
	"time"
	"github.com/prizarena/turn-based"
)

type PairsPlayerEntity struct {
	Created      time.Time             `datastore:"dc,noindex,omitempty"`
	LastMove     time.Time             `datastore:"dl,noindex,omitempty"`
	TurnsCount   int                   `datastore:"tc,noindex,omitempty"`
	MatchedCount int                   `datastore:"mc,noindex,omitempty"`
	MatchedItems string                `datastore:"mi,noindex,omitempty"`
	Open1        turnbased.CellAddress `datastore:"o1,noindex,omitempty"`
	Open2        turnbased.CellAddress `datastore:"o2,noindex,omitempty"`
}

const PairsPlayerKind = "P"

type PairsPlayer struct {
	db.StringID
	*PairsPlayerEntity
}

var _ db.EntityHolder = (*PairsPlayer)(nil)

func (PairsPlayer) Kind() string {
	return PairsPlayerKind
}

func (player PairsPlayer) Entity() interface{} {
	return player.PairsPlayerEntity
}

func (PairsPlayer) NewEntity() interface{} {
	return new(PairsPlayerEntity)
}

func (player *PairsPlayer) SetEntity(entity interface{}) {
	player.PairsPlayerEntity = entity.(*PairsPlayerEntity)
}
