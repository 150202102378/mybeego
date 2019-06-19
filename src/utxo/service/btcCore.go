package service

import (
	"os"
	"time"
	d "utxo/dao"
	"utxo/ini"
	m "utxo/models"
	u "utxo/utils"
)

//BTCCore obj
type BTCCore struct {
	StartBlock  int64 //数据开始记录的区块
	PointBlock  int64 //数据库当前最高的区块
	LatestBlock int64 //区块链最新的区块
}

//Init BTCCore
func (c *BTCCore) Init() {
	c.StartBlock = ini.GetStartBlock()
	//LastBlock
	c.iniLastBlock()
	//PointBlock
	c.iniPointBlock()
}

func (c *BTCCore) iniPointBlock() {
	ok, block := d.GetHighestBlock()
	if !ok || block.Height < c.StartBlock {
		c.PointBlock = c.StartBlock
		return
	}
	c.PointBlock = block.Height + 1
}

func (c *BTCCore) updatePointBlock() {
	ok, block := d.GetHighestBlock()
	if !ok {
		log.Fatal("DB ERROR: cannot get highest block")
		os.Exit(3)
	}
	c.PointBlock = block.Height + 1
}

func (c *BTCCore) iniLastBlock() { //name ？？？
	ok, result := CallRPC("getblockchaininfo", []string{})
	if !ok {
		log.Fatal("init BTCCore Error: cannot init LastBlock")
		os.Exit(3)
	}
	blocks := result.(map[string]interface{})["blocks"]
	c.LatestBlock = int64(blocks.(float64))
}

func getBlockInfo(height int64) (bool, map[string]interface{}) {
	result := make(map[string]interface{})
	if ok, hash := CallRPC("getblockhash", []int64{height}); ok {
		if ok, info := CallRPC("getblock", []string{hash.(string)}); ok {
			infoM := info.(map[string]interface{})
			result["hash"] = hash.(string)
			result["txs"] = infoM["tx"]
			result["confirmations"] = infoM["confirmations"]
			return true, result
		}
	}
	return false, result
}

func getTxidInfo(txid string) (bool, map[string]interface{}) {
	result := make(map[string]interface{})
	if ok, txidInfo := CallRPC("getrawtransaction", []interface{}{txid, true}); ok {
		if infoM, oko := txidInfo.(map[string]interface{}); oko {
			result["vin"] = infoM["vin"]
			result["vout"] = infoM["vout"]
			return true, result
		}
		return false, result
	}
	return false, result
}

func getVoutsAddresses(vouts []map[string]interface{}) []string {
	var addresses []string
	for _, vout := range vouts {
		address := vout["scriptPubKey"].(map[string]interface{})["addresses"]
		if address != nil {
			addresses = append(addresses, u.ConvertToSliceString(address.([]interface{}))...)
		}
	}
	return addresses
}

func getVinsAddresses(vins []map[string]interface{}) []string {
	var addresses []string
	for _, vin := range vins {
		if txid := vin["txid"]; txid != nil {
			if ok, txidInfo := getTxidInfo(txid.(string)); ok {
				vouts := txidInfo["vout"]
				address := getVoutsAddresses(u.ConvertToMapInterface(vouts.([]interface{})))
				addresses = append(addresses, address...)
			}
		}
	}
	return addresses
}

func revertTx(txid string) {
	ok, vouts := d.GetVoutsByTxid(txid)
	if ok {
		for _, vout := range vouts {
			delOk := d.DeleteVout(vout)
			if !delOk {
				log.Fatal("DB ERROR: cannot delete revertVout txid:", txid, " index :", vout.Index)
				os.Exit(3)
			}
		}
	}
	ok, _ = d.GetVinsByTxid(txid)
	if ok {
		delOK := d.DeleteVinByTxid(txid)
		if !delOK {
			log.Fatal("DB ERROR: cannot delete revertVinsByTxid ", txid)
			os.Exit(3)
		}
	}
	delTxOK := d.DeleteTx(txid)
	if !delTxOK {
		log.Fatal("DB ERROR: cannot delete revertTxs ", txid)
		os.Exit(3)
	}
}

func updateConfirm(updateNum int64) {
	//udpate blocks
	ok, blocks := d.GetUnconfirmedBlocks()
	if !ok {
		log.Fatal("DB ERROR: cannot get unconfirm blocks")
		os.Exit(3)
	} else {
		for _, block := range blocks {
			upOk := d.UpdateBlockConfirm(block, updateNum)
			if !upOk {
				log.Fatal("DB ERROR: cannot update block confirmations ", block.Height)
				os.Exit(3)
			}
		}
	}
	//update vins
	ok, vins := d.GetUnconfirmedVins()
	if !ok {
		log.Error("DB ERROR OR NONE VIN DATA")
	} else {
		for _, vin := range vins {
			upOk := d.UpdateVinConfirm(vin, updateNum)
			if !upOk {
				log.Fatal("DB ERROR: cannot update vin confirmations ", vin.Txid, " index:", *vin.Index)
				os.Exit(3)
			}
		}
	}
	//update vouts
	ok, vouts := d.GetUnconfirmedVouts()
	if !ok {
		log.Error("DB ERROR OR NONE VOUT DATA")
	} else {
		for _, vout := range vouts {
			upOk := d.UpdateVoutConfirm(vout, updateNum)
			if !upOk {
				log.Fatal("DB ERROR: cannot update vout confirmations ", vout.Txid, " index:", *vout.Index)
				os.Exit(3)
			}
		}
	}
	//update addr?
}

func checkXCurrencyTx(tx string, ch chan string, p *u.Pool) {
	ok, txInfo := getTxidInfo(tx)
	var isXCurrencyTx bool
	if !ok {
		log.Fatal("BTCCore ERROR: cannot get txInfo ", tx)
		os.Exit(3)
	}

	//check vout
	voutInfo := txInfo["vout"]
	voutAddress := getVoutsAddresses(u.ConvertToMapInterface(voutInfo.([]interface{})))
	if isXCurrencyAddress(voutAddress) {
		isXCurrencyTx = true
	}

	//check vin
	vinInfo := txInfo["vin"]
	if vinInfo != nil {
		vinAddress := getVinsAddresses(u.ConvertToMapInterface(vinInfo.([]interface{})))
		if isXCurrencyAddress(vinAddress) {
			isXCurrencyTx = true
		}
	}
	if isXCurrencyTx {
		ch <- tx
	}
	p.Done()
}

func filterXCurrencyTx(tallyBlockTxs []string) []string {
	var xCurrencyTxs []string
	var xCurrencyTxsCh = make(chan string, len(tallyBlockTxs))
	pool := u.NewPool(100)
	for _, tallyTx := range tallyBlockTxs {
		pool.Add(1)
		go checkXCurrencyTx(tallyTx, xCurrencyTxsCh, pool)
	}
	pool.Wait()
	for i := 0; i < len(xCurrencyTxsCh); i++ {
		x := <-xCurrencyTxsCh
		xCurrencyTxs = append(xCurrencyTxs, x)
	}
	return xCurrencyTxs
}

func tallyVout(xCurrencyTx string, fromAddress []string, datas []map[string]interface{}) {
	for _, xVout := range datas {
		xVoutAddrs := getVoutsAddresses([]map[string]interface{}{xVout})
		if isXCurrencyAddress(xVoutAddrs) {
			tnowStr := u.FormatTimeToDateTime(time.Now().UTC())
			var unspent int64
			unspent = 0
			var index int64
			index = int64(xVout["n"].(float64))
			newVout := m.Vout{
				Txid:          xCurrencyTx,
				Index:         &index, //sql.NullInt64{int64(xVout["n"].(float64)), true},
				Amount:        xVout["value"].(float64),
				Unspent:       &unspent,
				Confirmations: 1,
				FromAddr:      fromAddress,
				ToAddr:        u.ConvertToSliceString(xVout["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{}))[0],
				Created:       tnowStr,
				Updated:       tnowStr,
			}
			insertOk := d.UpdateVout(newVout)
			if !insertOk {
				log.Fatal("DB ERROR: cannot insert vout ", newVout.Txid, " index:", *newVout.Index)
				os.Exit(3)
			}
		}
	}
}

func tallyVin(xCurrencyTx string, toAddress []string, vinInfo interface{}) {
	if vinInfo == nil {
		return
	}
	for _, xVin := range u.ConvertToMapInterface(vinInfo.([]interface{})) {
		xVinAddress := getVinsAddresses([]map[string]interface{}{xVin})
		if isSysAddress(xVinAddress) {
			xVinTx := xVin["txid"]
			if xVinTx == nil {
				return
			}
			ok, xVinTxInfo := getTxidInfo(xVinTx.(string))
			if !ok {
				log.Fatal("BTCCore ERROR: cannot get txInfo ", xVinTx)
				os.Exit(3)
			}
			xVinTxVoutInfo := xVinTxInfo["vout"]
			for _, xVinTxVout := range u.ConvertToMapInterface(xVinTxVoutInfo.([]interface{})) {
				xVinVoutAddress := getVoutsAddresses([]map[string]interface{}{xVinTxVout})
				if isXCurrencyAddress(xVinVoutAddress) {
					tnowStr := u.FormatTimeToDateTime(time.Now().UTC())
					var index, spent int64
					index = int64(xVinTxVout["n"].(float64))
					spent = 0
					newVin := m.Vin{
						Txid:          xVinTx.(string), //xCurrencyTx,
						Index:         &index,          //sql.NullInt64{int64(xVinTxVout["n"].(float64)), true},
						Amount:        xVinTxVout["value"].(float64),
						Spend:         &spent,
						Confirmations: 1,
						FromAddr:      u.ConvertToSliceString(xVinTxVout["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})),
						ToAddr:        toAddress,
						Created:       tnowStr,
						Updated:       tnowStr,
					}
					insertOk := d.UpdateVin(newVin)
					if !insertOk {
						log.Fatal("DB ERROR: cannot insert vin ", newVin.Txid, " index:", *newVin.Index)
						os.Exit(3)
					}
				}
			}
		}
	}
}

func (c BTCCore) tallyTx(xCurrencyTxs []string, blockHash string) {
	if len(xCurrencyTxs) <= 0 {
		return
	}
	for _, xCurrencyTx := range xCurrencyTxs {
		ok, txInfo := getTxidInfo(xCurrencyTx)
		if !ok {
			log.Fatal("BTCCore ERROR: cannot get txInfo ", xCurrencyTx)
			os.Exit(3)
		}
		voutInfo := txInfo["vout"]
		vinInfo := txInfo["vin"]
		var fromAddress []string
		toAddress := getVoutsAddresses(u.ConvertToMapInterface(voutInfo.([]interface{})))
		if vinInfo != nil {
			fromAddress = getVinsAddresses(u.ConvertToMapInterface(vinInfo.([]interface{})))
		} else {
			fromAddress = []string{}
		}

		//tally vout
		tallyVout(xCurrencyTx, fromAddress, u.ConvertToMapInterface(voutInfo.([]interface{})))
		//tally vin
		tallyVin(xCurrencyTx, toAddress, vinInfo)
		//update tx table
		insertOk := d.UpdateTxs(m.Txs{
			Txid:      xCurrencyTx,
			BlockHash: blockHash,
		})
		if !insertOk {
			log.Fatal("DB ERROR: cannot insert txs ", xCurrencyTx)
			os.Exit(3)
		}
	}
}
