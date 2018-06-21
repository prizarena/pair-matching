package paircommands

import (
	"github.com/strongo/app"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
	"fmt"
	"strings"
)

func renderPairsBoardMessage(t strongo.SingleLocaleTranslator, tournament *pamodels.Tournament, board pairmodels.PairsBoard, players []pairmodels.PairsPlayer) (m bots.MessageFromBot, err error) {
	kbRows := make([][]tgbotapi.InlineKeyboardButton, board.SizeY)
	for y, row := range board.Rows() {
		if len(row) != board.SizeX {
			err = fmt.Errorf("len(board.Rows()[%v]) != board.SizeX: %v != %v", y, len(row), board.SizeX)
			return
		}
		kbRow := make([]tgbotapi.InlineKeyboardButton, board.SizeX)
		for x, cell := range row {
			var text string

			for _, player := range players {
				if player.OpenX == x && player.OpenY == y {
					text = string(cell)
					break
				} else if strings.Contains(player.MatchedItems, string(cell)) {
					text = "⬜"
					break
				}
			}
			if text == "" {
				text = "⬛"
			}
			kbRow[x] = tgbotapi.InlineKeyboardButton{Text: text, CallbackData: fmt.Sprintf("open?x=%v&y%v", x+1, y+1)}
		}
		kbRows[y] = kbRow
	}
	m.Text = "<b>Pair matching game</b>\nFind matching pairs as quickly as you can."
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
	return
}

