package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	d "utxo/dao"
	"utxo/ini"

	"github.com/fatih/set"
	"github.com/sirupsen/logrus"
)

var (
	log              *logrus.Logger
	httpClient       *http.Client
	rpcURL           string
	rpcUser          string
	rpcPasswd        string
	xCurrencyAddress set.Interface
	xSysAddress      set.Interface
)

func init() {
	log = ini.GetLog()
	httpClient = ini.GetHTTPClient()
	rpcUser, rpcPasswd, rpcURL = ini.GetRPCInfo()
	initXCurrencyAddress()
	initXSysAddress()
}

func initXCurrencyAddress() {
	var ok bool
	var xCurrencyAddr []string
	if ok, xCurrencyAddr = d.GetXCurrencyAddress(); !ok {
		log.Fatal("init XCurrency Addresses Error")
		os.Exit(3)
	}
	xCurrencyAddress = set.New(set.ThreadSafe)
	for _, xCurrency := range xCurrencyAddr {
		xCurrencyAddress.Add(xCurrency)
	}
}

func initXSysAddress() {
	xSysAddress = set.New(set.ThreadSafe)
	xSysAddress.Add("17SLoyVJUzMhkxBAo2gQHpsHD1qoo8h1b4")
}

//CallRPC call btc prc
func CallRPC(method string, params interface{}) (bool, interface{}) {
	result := make(map[string]interface{})
	data := map[string]interface{}{
		"method": method,
		"params": params,
	}
	dataStr, _ := json.Marshal(data)
	req, errNew := http.NewRequest("POST", rpcURL, bytes.NewReader(dataStr))
	req.SetBasicAuth(rpcUser, rpcPasswd)
	resp, errDo := httpClient.Do(req)
	if errNew != nil || errDo != nil {
		return false, result
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	if result["error"] != nil {
		return false, result
	}
	if result["result"] == nil {
		log.Fatal(string(body))
		return false, result
	}
	return true, result["result"]
}

func isXCurrencyAddress(addresses []string) bool {
	addressesSet := set.New(set.ThreadSafe)
	for _, addr := range addresses {
		addressesSet.Add(addr)
	}
	u := set.Intersection(xCurrencyAddress, addressesSet)
	if u.IsEmpty() {
		return false
	}
	return true
}

func isSysAddress(addresses []string) bool {
	addressesSet := set.New(set.ThreadSafe)
	for _, addr := range addresses {
		addressesSet.Add(addr)
	}
	u := set.Intersection(xSysAddress, addressesSet)
	if u.IsEmpty() {
		return false
	}
	return true
}
