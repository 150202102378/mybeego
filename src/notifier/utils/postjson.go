package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

//把json通过post请求发送出去，成功返回true，失败返回false
func Postjson(json []byte, url string) int {
	//fmt.Println(string(json))

	//构造request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json")
	//构造客户端并发送request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read failed")
	}
	return int(body[0] - '0')
}
