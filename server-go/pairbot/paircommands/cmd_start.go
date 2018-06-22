package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
	"net/url"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/prizarena/pair-matching/server-go/pairtrans"
	"fmt"
	"bytes"
)

const startCommandCommandCode = "start"

var startCommand = bots.Command{
	Code:     startCommandCommandCode,
	Commands: []string{"/start"},
	Action:   startAction,
	CallbackAction: func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		q := callbackUrl.Query()
		lang := q.Get("l")
		switch lang {
		case "ru":
			lang = "ru-RU"
		case "en":
			lang = "en-US"
		default:
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
				Text: "Unknown language",
			})
			return
		}
		if lang != "" {
			if lang == q.Get("cl") {
				m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
					Text: "It is already current language",
				})
				return
			} else {
				whc.SetLocale(lang)
				whc.ChatEntity().SetPreferredLanguage(lang)
			}
		}
		m, err = startAction(whc)
		m.IsEdit = true
		return
	},
}

func startAction(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {

	text := new(bytes.Buffer)
	fmt.Fprintf(text, "<b>%v</b>\n", whc.Translate(pairtrans.GameCardTitle))
	m.Text = text.String()
	m.Format = bots.MessageFormatHTML
	switchInlinePlay := "play"
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: "🇷🇺 Русский", CallbackData: "start?l=ru?&cl="+whc.Locale().Code5},
			{Text: "🇬🇧 English 🇺🇸", CallbackData: "start?l=en&cl="+whc.Locale().Code5},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: "⚔ Play", SwitchInlineQuery: &switchInlinePlay},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: "🏆 Tournaments", CallbackData: "tournaments"},
		},
	)
	return
}
