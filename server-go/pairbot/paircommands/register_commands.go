package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/prizarena-public/pabot"
	"net/http"
	"github.com/prizarena/prizarena-public/prizarena-interfaces"
	"github.com/prizarena/prizarena-public/prizarena-client-go"
	"github.com/prizarena/pair-matching/server-go/pairsecrets"
)

func RegisterPairCommands(router bots.WebhooksRouter) {
	router.RegisterCommands([]bots.Command{
		startCommand,
		inlineQueryCommand,
		openCellCommand,
		newBoardCommand,
	})

	pabot.InitPrizarenaBot("pairmatching", router, func(httpClient *http.Client) prizarena_interfaces.ApiClient {
		if httpClient == nil {
			panic("httpClient == nil")
		}
		return prizarena.NewHttpApiClient(httpClient, "", pairsecrets.PrizarenaGameID, pairsecrets.PrizarenaToken)
	})
}
