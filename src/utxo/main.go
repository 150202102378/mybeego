package main

import (
	"time"
	"utxo/ini"
	s "utxo/service"
)

var log = ini.GetLog()

func main() {
	log.Info("stating............")
	core := new(s.BTCCore)
	log.Info("new BTCCore........")
	core.Init()
	core.CheckRevert()
	log.Info("init BTCCore.......")
	for {
		if core.PointBlock < core.LatestBlock {
			log.Info("Tally Block : ", core.PointBlock)
			tallyTime := time.Now()
			core.Tally()
			log.Info("------------------------------ Tally Blocks Time : ", time.Since(tallyTime))
		} else {
			log.Info("Break Block : ", core.PointBlock)
			break
		}
	}
}
