package pairrouting

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/pair-matching/server-go/pairbot/paircommands"
)

var WebhooksRouter = bots.NewWebhookRouter(
	map[bots.WebhookInputType][]bots.Command{},
	func() string { return "Oops..." },
)

func init() {
	paircommands.RegisterPairCommands(WebhooksRouter)
}
