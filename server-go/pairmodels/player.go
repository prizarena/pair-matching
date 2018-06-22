package pairmodels

import (
	"github.com/strongo/db"
	"time"
	"github.com/prizarena/turn-based"
	"github.com/pkg/errors"
	"strings"
)

type PairsPlayerEntity struct {
	Created      time.Time             `datastore:"dc,noindex,omitempty"`
	LastMove     time.Time             `datastore:"dl,noindex,omitempty"`
	TurnsCount   int                   `datastore:"tc,noindex,omitempty"`
	UserName     string                `datastore:"un,noindex,omitempty"`
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

func (player PairsPlayer) UserID() string {
	if player.ID == "" {
		return ""
	}
	return player.ID[strings.Index(player.ID, ":")+1:]
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

func (player PairsPlayer) BeforeSave() error {
	counts := make(map[rune]int, player.MatchedCount)
	var matchedCount int
	for _, r := range player.MatchedItems {
		counts[r]++
		matchedCount++
	}
	for _, c := range counts {
		if c > 1 {
			return errors.New("duplicates in MatchedItems: " + player.MatchedItems)
		}
	}
	if player.MatchedCount != matchedCount {
		return errors.New("player.MatchedCount != count(player.MatchedItems)")
	}
	return nil
}
