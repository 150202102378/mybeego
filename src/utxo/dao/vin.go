package dao

import (
	"time"
	m "utxo/models"
	u "utxo/utils"
)

//GetVin get vin by txid and index
func GetVin(txid string, index int64) bool {
	var vin m.Vin
	err := psqlDB.Where("txid = ? and index = ?", txid, index).First(&vin).Error
	if err != nil {
		return false
	}
	return true
}

//UpdateVin update and insert vin data
func UpdateVin(vin m.Vin) bool {
	exist := GetVin(vin.Txid, *vin.Index)
	if exist {
		vin.Updated = u.FormatTimeToDateTime(time.Now().UTC())
		err := psqlDB.Model(&vin).Update("updated", vin.Updated).Error
		if err != nil {
			return false
		}
	} else {
		err := psqlDB.Create(&vin).Error
		if err != nil {
			log.Error("UpdateVin : ", err, " Vin :", vin)
			return false
		}
	}
	return true
}

//GetVinsByTxid get vins by txid
func GetVinsByTxid(txid string) (bool, []m.Vin) {
	var vins []m.Vin
	err := psqlDB.Where("txid = ?", txid).Find(&vins).Error
	if err != nil {
		return false, vins
	}
	return true, vins
}

//DeleteVinByTxid delete vin by txid
func DeleteVinByTxid(txid string) bool {
	err := psqlDB.Where("txid = ?", txid).Delete(&m.Vin{}).Error
	if err != nil {
		return false
	}
	return true
}

//GetUnconfirmedVins get all vins
func GetUnconfirmedVins() (bool, []m.Vin) {
	var vins []m.Vin
	err := psqlDB.Where("confirmations < 20").Find(&vins).Error
	if err != nil {
		return false, vins
	}
	return true, vins
}

//UpdateVinConfirm update vin confirm
func UpdateVinConfirm(vin m.Vin, updateNum int64) bool {
	vin.Confirmations += updateNum
	var err error
	if vin.Confirmations > confirmNum && *vin.Spend == 0 {
		*vin.Spend = 1
		err = psqlDB.Model(&vin).Update("spend", vin.Spend).Update("confirmations", vin.Confirmations).Error
	} else {
		err = psqlDB.Model(&vin).Update("confirmations", vin.Confirmations).Error
	}

	if err != nil {
		return false
	}
	return true
}
