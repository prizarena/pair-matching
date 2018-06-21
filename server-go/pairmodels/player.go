package pairmodels

import (
	"github.com/strongo/db"
	"time"
)

type PairsPlayerEntity struct {
	Created      time.Time `datastore:"dc,noindex,omitempty"`
	LastMove     time.Time `datastore:"dl,noindex,omitempty"`
	TurnsCount   int       `datastore:"tc,noindex,omitempty"`
	MatchedCount int       `datastore:"mc,noindex,omitempty"`
	MatchedItems string    `datastore:"mi,noindex,omitempty"`
	OpenX        int       `datastore:"ox,noindex,omitempty"`
	OpenY        int       `datastore:"oy,noindex,omitempty"`
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

func (player PairsPlayer) SetEntity(entity interface{}) {
	player.PairsPlayerEntity = entity.(*PairsPlayerEntity)
}
