package paircommands

import (
	"github.com/strongo/app"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
	"fmt"
	"strings"
	"github.com/prizarena/turn-based"
)

func renderPairsBoardMessage(t strongo.SingleLocaleTranslator, tournament *pamodels.Tournament, board pairmodels.PairsBoard, players []pairmodels.PairsPlayer) (m bots.MessageFromBot, err error) {
	width, height := board.Size.WidthHeight()
	kbRows := make([][]tgbotapi.InlineKeyboardButton, height)
	for y, row := range board.Rows() {
		if len(row) != width {
			err = fmt.Errorf("len(board.Rows()[%v]) != board.SizeX: %v != %v", y, len(row), width)
			return
		}
		kbRow := make([]tgbotapi.InlineKeyboardButton, width)
		for x, cell := range row {
			var text string

			for _, player := range players {
				if player.Open1.IsXY(x, y) || player.Open2.IsXY(x, y) {
					text = string(cell)
					break
				} else if strings.Contains(player.MatchedItems, string(cell)) {
					text = " "
					break
				}
			}
			if text == "" {
				text = "⬜"
			}
			kbRow[x] = tgbotapi.InlineKeyboardButton{Text: text, CallbackData: fmt.Sprintf("open?c=%v&b=%v", turnbased.NewCellAddress(x, y), board.ID)}
		}
		kbRows[y] = kbRow
	}
	m.Text = "<b>Pair matching game</b>\nFind matching pairs as quickly as you can."
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
	return
}

