package dao

import (
	m "utxo/models"
)

//UpdateTrafTask update and insert data
func UpdateTrafTask(trafTask m.TransferTasks) bool {
	err := psqlDB.Create(trafTask).Error
	if err != nil {
		return false
	}
	return true
}
