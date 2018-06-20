package pairgaedal

import (
	"github.com/strongo/db/gaedb"
	"github.com/prizarena/pair-matching/server-go/pairdal"
)

func RegisterDal() {
	pairdal.DB = gaedb.NewDatabase()
}
