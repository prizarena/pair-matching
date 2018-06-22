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
	"bytes"
)

func renderPairsBoardMessage(t strongo.SingleLocaleTranslator, tournament *pamodels.Tournament, board pairmodels.PairsBoard, players []pairmodels.PairsPlayer) (m bots.MessageFromBot, err error) {
	isCompleted := board.IsCompleted(players)
	m.Format = bots.MessageFormatHTML
	text := new(bytes.Buffer)
	fmt.Fprintln(text, `<a href="https://t.me/PairMatchingGameBot">Pair matching game</a>`)
	fmt.Fprintln(text, "Find matching pairs as quickly as you can.")
	for i, p := range players {
		fmt.Fprintf(text, "%d. <b>%v</b>: %v\n", i+1, p.UserName, p.MatchedCount)
	}
	if isCompleted {
		text.WriteString("\n<b>Board:</b>")
		text.WriteString(board.DrawBoard())
		text.WriteString("\n<b>Choose size of next board:</b>")
		m.Keyboard = newBoardSizesKeyboard
	} else {
		width, height := board.Size.WidthHeight()
		kbRows := make([][]tgbotapi.InlineKeyboardButton, height)
		for y, row := range board.Rows() {
			if len(row) != width {
				err = fmt.Errorf("len(board.Rows()[%v]) != board.SizeX: %v != %v", y, len(row), width)
				return
			}
			kbRow := make([]tgbotapi.InlineKeyboardButton, width)
			const (
				isMatched = " "
				closed = "⬜"
			)
			for x, cell := range row {
				var text string

				for _, player := range players {
					if player.Open1.IsXY(x, y) || player.Open2.IsXY(x, y) {
						text = string(cell)
						break
					} else if strings.Contains(player.MatchedItems, string(cell)) {
						text = isMatched
						break
					}
				}
				if text == "" {
					text = closed
				}
				kbRow[x] = tgbotapi.InlineKeyboardButton{Text: text, CallbackData: openCellCallbackData(turnbased.NewCellAddress(x, y), board.ID)}
			}
			kbRows[y] = kbRow
		}
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
	}
	m.Text = text.String()
	return
}

