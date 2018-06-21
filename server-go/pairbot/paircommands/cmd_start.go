package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
)

const startCommandCommandCode = "start"

var startCommand = bots.Command{
	Code: startCommandCommandCode,
	Commands: []string{"/start"},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		m.Text = "<b>Pair-Matching game</b>"
		m.Format = bots.MessageFormatHTML
		switchInlinePlay := "play"
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: "Play", SwitchInlineQuery: &switchInlinePlay},
			},
		)
		return
	},
}
