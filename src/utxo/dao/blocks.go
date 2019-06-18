package dao

import (
	m "utxo/models"
)

//UpdateBlock new blocks data
func UpdateBlock(block m.Blocks) bool {
	err := psqlDB.Create(&block).Error
	if err != nil {
		return false
	}
	return true
}

//GetHighestBlock get highest block
func GetHighestBlock() (bool, m.Blocks) {
	var block m.Blocks
	err := psqlDB.Order("height desc").First(&block).Error
	if err != nil {
		return false, block
	}
	return true, block
}

//GetBlockByHeight get block by height
func GetBlockByHeight(height int64) (bool, m.Blocks) {
	var block m.Blocks
	err := psqlDB.Where("height = ?", height).First(&block).Error
	if err != nil {
		return false, block
	}
	return true, block
}

//DeleteBlockByHeight delete block by height
func DeleteBlockByHeight(height int64) bool {
	err := psqlDB.Where("height = ?", height).Delete(&m.Blocks{}).Error
	if err != nil {
		return false
	}
	return true
}

//GetUnconfirmedBlocks get all blocks
func GetUnconfirmedBlocks() (bool, []m.Blocks) {
	var blocks []m.Blocks
	err := psqlDB.Where("confirmations < 20").Find(&blocks).Error
	if err != nil {
		return false, blocks
	}
	return true, blocks
}

//UpdateBlockConfirm update block confirmations
func UpdateBlockConfirm(block m.Blocks, updateNum int64) bool {
	block.Confirmations += updateNum
	var err error
	if block.Confirmations > confirmNum {
		block.Confirm = true
		err = psqlDB.Model(&block).Update("confirmations", block.Confirmations).Update("confirm", block.Confirm).Error
	} else {
		err = psqlDB.Model(&block).Update("confirmations", block.Confirmations).Error
	}
	if err != nil {
		return false
	}
	return true
}
