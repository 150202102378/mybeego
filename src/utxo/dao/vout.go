package dao

import (
	"time"
	m "utxo/models"
	u "utxo/utils"
)

//GetVout get vout by txid and index
func GetVout(txid string, index int64) bool {
	var vout m.Vout
	err := psqlDB.Where("txid = ? and index = ?", txid, index).First(&vout).Error
	if err != nil {
		return false
	}
	return true
}

//UpdateVout update and insert vout data
func UpdateVout(vout m.Vout) bool {
	exist := GetVout(vout.Txid, *vout.Index)
	if exist {
		vout.Updated = u.FormatTimeToDateTime(time.Now().UTC())
		err := psqlDB.Model(&vout).Update("updated", vout.Updated).Error
		if err != nil {
			return false
		}
	} else {
		err := psqlDB.Create(&vout).Error
		ok := UpdateAddr(vout.ToAddr, vout.Amount)
		if err != nil || !ok {
			return false
		}
	}
	return true
}

//GetVoutsByTxid get vouts by txid
func GetVoutsByTxid(txid string) (bool, []m.Vout) {
	var vouts []m.Vout
	err := psqlDB.Where("txid = ?", txid).Find(&vouts).Error
	if err != nil {
		return false, vouts
	}
	return true, vouts
}

//GetVoutsByToAddr get vouts by toAddr
func GetVoutsByToAddr(toAddr string) (bool, []m.Vout) {
	var vouts []m.Vout
	err := psqlDB.Where("toAddr = ?", toAddr).Find(&vouts).Error
	if err != nil {
		return false, vouts
	}
	return true, vouts
}

//DeleteVout delete
func DeleteVout(vout m.Vout) bool {
	ok, addr := GetAddr(vout.Txid)
	if !ok {
		return false
	}
	UpdateAddrConfirm(addr, -vout.Amount, false)
	err := psqlDB.Delete(&vout).Error
	if err != nil {
		return false
	}
	return true
}

//DeleteVoutsByTxid delete vouts by txid
func DeleteVoutsByTxid(txid string) bool {
	err := psqlDB.Where("txid = ?", txid).Delete(&m.Vout{}).Error
	if err != nil {
		return false
	}
	return true
}

//GetUnconfirmedVouts get all vouts
func GetUnconfirmedVouts() (bool, []m.Vout) {
	var vouts []m.Vout
	err := psqlDB.Where("confirmations < 20").Find(&vouts).Error
	if err != nil {
		return false, vouts
	}
	return true, vouts
}

//UpdateVoutConfirm update vout confirm
func UpdateVoutConfirm(vout m.Vout, updateNum int64) bool {
	vout.Confirmations += updateNum
	var err error
	if vout.Confirmations > confirmNum && *vout.Unspent == 0 {
		*vout.Unspent = 1
		err = psqlDB.Model(&vout).Update("unspent", vout.Unspent).Update("confirmations", vout.Confirmations).Error
		ok, addr := GetAddr(vout.ToAddr)
		if !ok {
			return false
		}
		ok = UpdateAddrConfirm(addr, vout.Amount, true)
		if !ok {
			return false
		}
	} else {
		err = psqlDB.Model(&vout).Update("confirmations", vout.Confirmations).Error
	}
	if err != nil {
		return false
	}
	return true
}
