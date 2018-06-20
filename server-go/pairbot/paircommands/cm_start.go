package paircommands

import "github.com/strongo/bots-framework/core"

const startCommandCommandCode = "start"

var startCommand = bots.Command{
	Code: startCommandCommandCode,
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		return
	},
}