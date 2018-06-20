package pairapp

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsbot"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/prizarena/pair-matching/server-go/pairdal/pairgaedal"
)

func InitApp(botHost bots.BotHost) {
	pairgaedal.RegisterDal()

	httpRouter := httprouter.New()
	http.Handle("/", httpRouter)

	rpsbot.InitBot(httpRouter, botHost, pairAppContext{})
}
