package pairmodels

import (
	"github.com/strongo/db"
	"bytes"
	"github.com/prizarena/turn-based"
	"math/rand"
	"time"
	"fmt"
)

type PairsBoardEntity struct {
	Cells string `datastore:",noindex,omitempty"`
	SizeX int    `datastore:",noindex"`
	SizeY int    `datastore:",noindex"`
	turnbased.BoardEntityBase
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

func (eh *PairsBoard) SetEntity(entity interface{}) {
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

func (board PairsBoardEntity) GetCell(x, y int) rune {
	if x <= 0 {
		panic(fmt.Sprintf("x <= 0: %v", x))
	}
	if y <= 0 {
		panic(fmt.Sprintf("y <= 0: %v", y))
	}
	x--
	y--
	k := y * board.SizeX + x
	var runeIndex int
	for _, r := range board.Cells {
		if runeIndex == k {
			return r
		}
		runeIndex++
	}
	return 0
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

func emojiSet() []rune {
	return []rune{
		'🚀',
		'🚁',
		'🚂',
		'🚃',
		'🚄',
		'🚅',
		'🚆',
		'🚇',
		'🚈',
		'🚉',
		'🚊',
		'🚋',
		'🚌',
		'🚍',
		'🚎',
		'🚏',
		'🚐',
		'🚑',
		'🚒',
		'🚓',
		'🚔',
		'🚕',
		'🚖',
		'🚗',
		'🚘',
		'🚙',
		'🚚',
		'🚛',
		'🚜',
		'🚝',
		'🚞',
		'🚟',
		'🚠',
		'🚡',
		'🚢',
		'🚣',
		'🚤',
		'🚥',
		'🚦',
		'🚧',
		'🚨',
		'🚩',
		'🚪',
		'🚫',
		'🚬',
		'🚭',
		'🚮',
		'🚯',
		'🚰',
		'🚱',
		'🚲',
		'🚳',
		'🚴',
		'🚵',
		'🚶',
		'🚷',
		'🚸',
		'🚹',
		'🚺',
		'🚻',
		'🚼',
		'🚽',
		'🚾',
		'🚿',
	}
}

func Shuffle(width, height int) string {
	available := emojiSet()
	pairsCount := width * height / 2

	items := make([]rune, pairsCount*2)
	for i := 0; i < pairsCount; i++ {
		randIndex := rand.Intn(len(available))
		items[i] = available[randIndex]
		items[i+pairsCount] = available[randIndex]
		available = append(available[:randIndex], available[randIndex+1:]...)
	}
	shuffle(rand.New(rand.NewSource(time.Now().UnixNano())), len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	s := new(bytes.Buffer)
	for _, r := range items {
		s.WriteRune(r)
	}
	return s.String()
}

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func shuffle(r *rand.Rand, n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(r.Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(r.Int31n(int32(i + 1)))
		swap(i, j)
	}
}
