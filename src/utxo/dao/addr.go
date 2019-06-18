package dao

import (
	m "utxo/models"
)

//NewAddr new addr data
func NewAddr(ad m.Addr) bool {
	err := psqlDB.Create(&ad).Error
	if err != nil {
		return false
	}
	return true
}

//GetAddr get addr by id
func GetAddr(id string) (bool, m.Addr) {
	var addr m.Addr
	err := psqlDB.Where("id = ?", id).Find(&addr).Error
	if err != nil {
		return false, addr
	}
	return true, addr
}

//GetAddrByConfirm get addr by confirm
func GetAddrByConfirm(confirm float64) (bool, []m.Addr) {
	var addrs []m.Addr
	err := psqlDB.Where("confirm > ?", confirm).Find(&addrs)
	if err != nil {
		return false, addrs
	}
	return true, addrs
}

//UpdateAddrConfirm update addr confirm amount
func UpdateAddrConfirm(addr m.Addr, amount float64, confirmed bool) bool {
	if confirmed {
		addr.Unconfirm -= amount
		addr.Confirm += amount
		err := psqlDB.Model(&addr).Update("unconfirm", addr.Unconfirm).
			Update("confirm", addr.Confirm).Error
		if err != nil {
			return false
		}
		return true
	}
	addr.Unconfirm += amount
	err := psqlDB.Model(&addr).Update("unconfirm", addr.Unconfirm).Error
	if err != nil {
		return false
	}
	return true
}

//UpdateAddr update addr data
func UpdateAddr(id string, amount float64) bool {
	var addr m.Addr
	err := psqlDB.Where("id = ?", id).Find(&addr).Error
	if err != nil {
		addr = m.Addr{
			ID:        id,
			Confirm:   0,
			Unconfirm: amount,
		}
		ok := NewAddr(addr)
		if !ok {
			return false
		}
		return true
	}
	return UpdateAddrConfirm(addr, amount, false)
}
