package models

import (
	_ "github.com/lib/pq"
)

type ConfirmedTx struct {
	Txid    string `orm:"pk"`
	Status  int
	Updated string
	To_addr string
	Index   int
	Amount  float64
}

type Draw struct {
	Id      string `orm:"pk"`
	Status  int
	Updated string
	Amount  float64
	To_addr string
}

type Json struct {
	Txid    string `json:"txid"`
	Index   int
	To_addr string
	Amount  float64
}
