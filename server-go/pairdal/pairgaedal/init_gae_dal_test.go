package pairgaedal

import (
	"testing"
	"github.com/prizarena/pair-matching/server-go/pairdal"
)

func TestRegisterDal(t *testing.T) {
	RegisterDal()
	if pairdal.DB == nil {
		t.Fatal("pairdal.DB == nil")
	}
}
