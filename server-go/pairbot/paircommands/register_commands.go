package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/prizarena-public/pabot"
				"github.com/prizarena/pair-matching/server-go/pairsecrets"
)

func RegisterPairCommands(router bots.WebhooksRouter) {
	router.RegisterCommands([]bots.Command{
		startCommand,
		inlineQueryCommand,
		openCellCommand,
		newBoardCommand,
		newPlayCommand,
	})

	pabot.InitPrizarenaInGameBot(pairsecrets.PrizarenaGameID, pairsecrets.PrizarenaToken, router)
}
