package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/pair-matching/server-go/pairsecrets"
)

const startCommandCommandCode = "start"

var startCommand = bots.Command{
	Code: startCommandCommandCode,
	Commands: []string{"/start"},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		if m, err = pabot.OnStartIfTournamentLink(whc, pairsecrets.PrizarenaGameID); m.Text != "" || err != nil {
			return
		}
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
