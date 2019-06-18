package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

//Blocks : struct of blocks
type Blocks struct {
	Height        int64  `gorm:"type:int;primary_key"`
	Hash          string `gorm:"type:text"`
	Confirmations int64
	Confirm       bool
	DealtTime     string
}

//Txs is struct of transaction
type Txs struct {
	Txid      string `gorm:"type:text;primary_key"`
	BlockHash string `gorm:"type:text"`
}

//Vout is struct of transactions
type Vout struct {
	Txid          string `gorm:"type:text;primary_key"`
	Index         *int64 `gorm:"type:SMALLINT;primary_key;not null;default:0"`
	Amount        float64
	Unspent       *int64 `gorm:"default:0"` // 0 : unspent;  1 : confirmed; 2 : transferred to sys_addr
	Confirmations int64
	FromAddr      pq.StringArray `gorm:"type:text[]"`
	ToAddr        string         `gorm:"type:text"`
	Created       string
	Updated       string
}

//Vin is struct of transactions
type Vin struct {
	Txid          string `gorm:"type:text;primary_key"`
	Index         *int64 `gorm:"type:SMALLINT;primary_key;default:0"`
	Amount        float64
	Spend         *int64 `gorm:"default:0"` // 0 : unconfirmed;  1 : spented and confirmed
	Confirmations int64
	FromAddr      pq.StringArray `gorm:"type:text[]"`
	ToAddr        pq.StringArray `gorm:"type:text[]"`
	Created       string
	Updated       string
}

//Addr is struct of transactions
type Addr struct {
	ID        string `gorm:"type:text;primary_key"`
	Confirm   float64
	Unconfirm float64
}

//XcurrencyAddress is struct of our xcurrency address
type XcurrencyAddress struct {
	Addr   string `gorm:"type:text;primary_key"`
	Userid string
}

//TransferTasks struct of transfer to sysAddr tasks
type TransferTasks struct {
	ID      string `gorm:"type:text;primary_key"`
	Status  *int64
	Sources postgres.Jsonb `gorm:"type:json"`
}

//Source struct of transfer source info
type Source struct {
	Txid    string
	Amount  float64
	Address string
	Index   *int64
}
