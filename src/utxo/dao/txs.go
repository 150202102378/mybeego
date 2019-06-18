package dao

import (
	m "utxo/models"
)

//UpdateTxs update and insert txs data
func UpdateTxs(tx m.Txs) bool {
	err := psqlDB.Create(&tx).Error
	if err != nil {
		return false
	}
	return true
}

//GetTxsByBlockHash get txs data by block hash
func GetTxsByBlockHash(blockHash string) (bool, []m.Txs) {
	var txs []m.Txs
	err := psqlDB.Where("block_hash = ?", blockHash).Find(&txs).Error
	if err != nil {
		return false, txs
	}
	return true, txs
}

//DeleteTx delete tx by txid
func DeleteTx(txid string) bool {
	err := psqlDB.Where("txid = ?", txid).Delete(&m.Txs{}).Error
	if err != nil {
		return false
	}
	return true
}
