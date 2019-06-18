package dao

import (
	"os"
	"utxo/ini"
	m "utxo/models"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	psqlDB     *gorm.DB
	log        *logrus.Logger
	confirmNum int64
)

func init() {
	log = ini.GetLog()
	confirmNum = ini.GetConfirmNum()
	initTable(m.Blocks{})
	initTable(m.Txs{})
	initTable(m.Vout{})
	initTable(m.Vin{})
	initTable(m.Addr{})
	initTable(m.TransferTasks{})
	//Tmp init
	initTable(m.XcurrencyAddress{})
}

func initTable(structData interface{}) {
	if psqlDB == nil {
		psqlDB = ini.GetPsqlDB()
	}
	if !psqlDB.HasTable(structData) {
		if err := psqlDB.CreateTable(structData).Error; err != nil {
			log.Fatal("init table error")
			os.Exit(3)
		}
	}
}
