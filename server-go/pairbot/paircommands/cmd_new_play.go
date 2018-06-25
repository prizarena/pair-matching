package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"github.com/strongo/bots-api-telegram"
	"fmt"
	"strconv"
	"github.com/strongo/app"
	"bytes"
	"github.com/prizarena/pair-matching/server-go/pairtrans"
	"github.com/prizarena/prizarena-public/pamodels"
)

const newPlayCommandCode = "new_play"

var newPlayCommand = bots.Command{
	Code: newPlayCommandCode,
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		return newPlayAction(whc, "")
	},
	CallbackAction: func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		tournamentID := callbackUrl.Query().Get("t")
		return newPlayAction(whc, tournamentID)
	},
}

func newPlayAction(whc bots.WebhookContext, tournamentID string) (m bots.MessageFromBot, err error) {
	var tournament pamodels.Tournament
	m.Text = getNewPlayText(whc, tournament)
	m.Format = bots.MessageFormatHTML
	m.Keyboard = getNewPlayTgInlineKbMarkup(whc.Locale().Code5, tournamentID)
	return
}

func getNewPlayTgInlineKbMarkup(lang, tournamentID string) *tgbotapi.InlineKeyboardMarkup {
	sizeButton := func(width, height int) tgbotapi.InlineKeyboardButton {
		return tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(strconv.Itoa(width) + "x" + strconv.Itoa(height)),
			CallbackData: getNewBoardCallbackData(width, height, tournamentID, lang),
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(4, 2),
			sizeButton(4, 3),
			sizeButton(4, 4),
			sizeButton(5, 4),
		},
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(6, 4),
			sizeButton(6, 5),
			sizeButton(6, 6),
		},
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(7, 6),
			sizeButton(8, 6),
			sizeButton(8, 7),
			sizeButton(8, 8),
		},
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(8, 9),
			sizeButton(8, 10),
			sizeButton(8, 11),
			sizeButton(8, 12),
		},
	)
}

var newNonTournamentBoardSizesKeyboards = map[string]*tgbotapi.InlineKeyboardMarkup{
	"en-US": getNewPlayTgInlineKbMarkup("en-US", ""),
	"ru-RU": getNewPlayTgInlineKbMarkup("ru-RU", ""),
}


func getNewPlayText(t strongo.SingleLocaleTranslator, tournament pamodels.Tournament) string {
	text := new(bytes.Buffer)
	text.WriteString(t.Translate(pairtrans.NewGameText))
	return text.String()
}