package logic

import (
	"encoding/json"
	"fmt"
	"notifier/dao"
	"notifier/utils"
	"sync"
)

/* //通知器，发送一个json数组切片
func otifier() bool {
	//获取status为0的数据，confirms是[]model.ConfirmedTx
	confirms := dao.RechargeQuery()

	//转换成json，只保留Txid，Index，To_addr，Amount这四个字段
	var m1 map[string]interface{}
	var m2 []map[string]interface{}
	for _, v := range confirms {
		m1 = make(map[string]interface{})
		m1["Txid"] = v.Txid
		m1["Index"] = v.Index
		m1["To_addr"] = v.To_addr
		m1["Amount"] = v.Amount
		m2 = append(m2, m1)
	}
	m3 := make(map[string]interface{})
	m3["datas"] = m2

	js, err := json.Marshal(m3)
	if err != nil {
		return false
	}

	fmt.Println(string(js))
	//输出结果：{"datas":[{"Amount":124.231,"Index":10,"To_addr":"312312","Txid":"3123123"},{"Amount":124.231,"Index":10,"To_addr":"312312","Txid":"3123124"},{"Amount":124.231,"Index":10,"To_addr":"312312","Txid":"3123125"}]}
	//把json通过post方法发送
	//utils.Postjson(js)
	return true
} */

func recharge(wg sync.WaitGroup) {
	//获取status为0的数据，confirms是[]model.ConfirmedTx
	confirms := dao.RechargeQuery()
	temp := make(map[string]interface{})
	for _, v := range confirms {
		//转换成json，只保留Txid，Index，To_addr，Amount这四个字段
		//fmt.Println(v)
		temp["txid"] = v.Txid
		temp["index"] = v.Index
		temp["to_addr"] = v.To_addr
		temp["amount"] = v.Amount
		js, err := json.Marshal(temp)
		fmt.Println("confirms json:", string(js))
		if err != nil {
			return
		}

		//把json通过post方法发送，并根据返回值修改confirm_tx表
		status := utils.Postjson(js, "http://192.168.1.46:8217/api/charge")
		fmt.Println("confirms json返回值：", status)
		dao.RechargeUpdate(v, status)
	}
	wg.Done()
}

func withdraw(wg sync.WaitGroup) {
	//获取status为0的数据，draws是[]model.Draw
	draws := dao.WithdrawQuery()
	//fmt.Println(draws)
	temp := make(map[string]interface{})
	for _, v := range draws {
		//检验地址是否正确,错误则修改status为4
		/* if !utils.Verify(v.To_addr, "Testnet") {

		} */

		//转换json，只保留id，amount，to_addr字段
		temp["id"] = v.Id
		temp["amount"] = v.Amount
		temp["to_addr"] = v.To_addr
		js, err := json.Marshal(temp)
		fmt.Println("draws json:", string(js))
		if err != nil {
			return
		}

		//把json通过post方法发送，并根据返回值修改draw表
		status := utils.Postjson(js, "http://192.168.1.46:8217/api/draw")
		fmt.Println("draws json返回值：", status)
		dao.WithdrawUpdate(v, status)
	}
	wg.Done()
}

/*
 *通知器
 *Notifier()是给外面调用的接口
 *功能：调用后会获取draw表和confrim_tx表status为0的信息，
 *     封装成json，然后post给相对应的url，根据respon信息来修改表的status
 *在提币过程中还需要给出一个接口去获取post请求，这个交给router来操作
 */
func Notifier() {
	//防止主线程在协程完成前关闭
	var wg sync.WaitGroup
	wg.Add(2)
	//充币所需要的通知操作
	go recharge(wg)
	//提币所需要的通知操作
	go withdraw(wg)
}
