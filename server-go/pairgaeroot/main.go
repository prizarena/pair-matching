package pairgaeroot

import (
	"github.com/strongo/log"
	"github.com/strongo/bots-framework/hosts/appengine"
	"github.com/prizarena/pair-matching/server-go/pairapp"
)

func init() {
	log.AddLogger(gaehost.GaeLogger)
	pairapp.InitApp(gaehost.GaeBotHost{})
}
