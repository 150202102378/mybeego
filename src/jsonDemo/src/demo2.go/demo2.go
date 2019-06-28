package demo2

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
