package service

import (
	"encoding/json"
	"os"
	"time"
	"utxo/dao"
	d "utxo/dao"
	m "utxo/models"
	"utxo/utils"
	u "utxo/utils"

	"github.com/jinzhu/gorm/dialects/postgres"
)

//CheckRevert to check BlockChain revert
func (c BTCCore) CheckRevert() {
	if c.PointBlock > c.LatestBlock {
		log.Debug("Revert......")
		c.Revert()
	}
}

//Revert Executor
func (c BTCCore) Revert() {
	revertNum := c.PointBlock - c.LatestBlock
	for i := 1; int64(i) < revertNum; i++ {
		revertHeight := c.LatestBlock + int64(i)
		ok, block := d.GetBlockByHeight(revertHeight)
		if !ok {
			log.Fatal("DB ERROR: cannot find revertBlock ", revertHeight)
			os.Exit(3)
		}
		ok, revertTxs := d.GetTxsByBlockHash(block.Hash)
		if ok {
			for _, tx := range revertTxs {
				revertTx(tx.Txid)
			}
		}
		delBlockOk := d.DeleteBlockByHeight(revertHeight)
		if !delBlockOk {
			log.Fatal("DB ERROR: cannot delete revertBlock ", revertHeight)
			os.Exit(3)
		}
	}
	//Update Confirm
	updateConfirm(-revertNum)
}

//Tally execute tally new block
func (c *BTCCore) Tally() {
	updateConfirm(1)
	tallyBlockHeight := c.PointBlock
	ok, tallyBlockInfo := getBlockInfo(tallyBlockHeight)
	if !ok {
		log.Fatal("BTCCore ERROR : cannot get tallyBlockInfo ", tallyBlockHeight)
		os.Exit(3)
	}
	xCurrencyTxs := filterXCurrencyTx(u.ConvertToSliceString(tallyBlockInfo["txs"].([]interface{})))
	//Tally xCurrencyTxs
	c.tallyTx(xCurrencyTxs, tallyBlockInfo["hash"].(string))
	d.UpdateBlock(m.Blocks{
		Height:        tallyBlockHeight,
		Hash:          tallyBlockInfo["hash"].(string),
		Confirmations: int64(tallyBlockInfo["confirmations"].(float64)),
		Confirm:       false,
		DealtTime:     u.FormatTimeToDateTime(time.Now().UTC()),
	})
	c.updatePointBlock()
}

//Transfer recording transfer
//Transfer应该加入到build上
func (c *BTCCore) Transfer() {
	//get all the addr that confirms > 0.8
	ok, addrs := dao.GetAddrByConfirm(0.8)
	if !ok {
		log.Fatal("BTCCore ERROR : cannot get addrs")
		os.Exit(3)
	}

	for _, addr := range addrs {
		//get vout
		ok, vouts := dao.GetVoutsByToAddr(addr.ID)
		if !ok {
			log.Fatal("BTCCore ERROR : cannot get vout")
			os.Exit(3)
		}

		//make transfer
		var t m.TransferTasks
		t.ID = utils.NewUUID().String()
		status := int64(0)
		t.Status = &status
		sources := make([]m.Source, len(vouts))
		for i, vout := range vouts {
			sources[i].Txid = vout.Txid
			sources[i].Amount = vout.Amount
			sources[i].Index = vout.Index
			sources[i].Address = vout.ToAddr
			//update vout, change the Unspent to 2
			unspent := int64(2)
			vout.Unspent = &unspent
			dao.UpdateVout(vout)
		}
		js, err := json.Marshal(sources)
		if err != nil {
			log.Fatal("BTCCore ERROR : make json faild", js)
		}
		t.Sources = postgres.Jsonb{RawMessage: js}

		//insert into transfer_task
		ok2 := dao.UpdateTrafTask(t)
		if !ok2 {
			log.Fatal("BTCCore ERROR : insert into transfer_tasks failed", t)
		}
	}

}
