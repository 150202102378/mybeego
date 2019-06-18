package dao

import (
	m "utxo/models"
)

//GetXCurrencyAddress : get xcurrency all address
func GetXCurrencyAddress() (bool, []string) {
	var address []string
	if err := psqlDB.Model(&m.XcurrencyAddress{}).Pluck("addr", &address).Error; err != nil {
		return false, address
	}
	return true, address
}
