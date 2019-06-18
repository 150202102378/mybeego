package logic

import (
	"encoding/json"
	"fmt"
	"notifier/dao"
	"notifier/models"
)

/* //根据post请求的内容更新draw表
func update(w http.ResponseWriter, r *http.Request) {
	//读取post请求的内容
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read fail")
		return
	}

	//反序列化，问题：json有什么字段？
	var dr models.Draw
	err1 := json.Unmarshal(body, &dr)
	if err1 != nil {
		fmt.Println("unmarshal fail")
		return
	}

	//更新drew
	dao.WithdrawUpdate(&dr, 2)
}

//给出接口：本机ip/update:12555 接收post请求
func startListen() {

	http.HandleFunc("/update/post", update)
	err := http.ListenAndServe(":12555", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
*/

//处理传过来的json，并更新draw表
func HandlePost(body []byte) {
	//反序列化，问题：json有什么字段？
	var dr models.Draw
	err := json.Unmarshal(body, &dr)
	if err != nil {
		fmt.Println("unmarshal fail")
		return
	}

	//更新draw
	dao.WithdrawUpdate(&dr, 2)
}
