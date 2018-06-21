package pairapp

import (
	"github.com/strongo/bots-framework/core"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/prizarena/pair-matching/server-go/pairdal/pairgaedal"
	"github.com/prizarena/pair-matching/server-go/pairbot"
)

func InitApp(botHost bots.BotHost) {
	pairgaedal.RegisterDal()

	httpRouter := httprouter.New()
	http.Handle("/", httpRouter)

	pairbot.InitBot(httpRouter, botHost, pairAppContext{})
}
