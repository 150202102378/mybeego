package service

import (
	"testing"
)

func TestInitBTCCore(t *testing.T) {
	core := new(BTCCore)
	core.Init()
}
