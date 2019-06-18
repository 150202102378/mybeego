package ini

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	//:"time"
	//init postgresql driver
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//PsqlConf is struct of psql
type PsqlConf struct {
	Host, Port, DB, User, Passwd string
}

//EmailConf is struct of email
type EmailConf struct {
	SMTPServer, SMTPPasswd, SenderEmail string
	ReceiveEmails                       []string
}

//BTCCoreConf is struct of btcconfig
type BTCCoreConf struct {
	RPCHost    string
	RPCPort    string
	RPCUser    string
	RPCPasswd  string
	ConfirmNum int64
	StartBlock int64
}

func (c BTCCoreConf) getURL() string {
	return fmt.Sprintf(
		"http://%s:%s/",
		c.RPCHost,
		c.RPCPort,
	)
}

//GetPsqlDB : get and init psql db
func GetPsqlDB() *gorm.DB {
	if psqlDB != nil {
		return psqlDB
	}
	psqlDBLock.Lock()
	defer psqlDBLock.Unlock()
	if psqlDB != nil {
		return psqlDB
	}
	db := config.Psql[runMode].DB
	host := config.Psql[runMode].Host
	port := config.Psql[runMode].Port
	user := config.Psql[runMode].User
	passwd := config.Psql[runMode].Passwd
	dbsql := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, db,
	)
	var err error
	psqlDB, err = gorm.Open("postgres", dbsql)
	if err != nil {
		log.Fatal("init psql err : ", err)
		os.Exit(3)
	}
	return psqlDB
}

//GetHTTPClient get and init http client
func GetHTTPClient() *http.Client {
	if httpClient == nil {
		//tr := &http.Transport{
		//	//MaxIdleConns:       200,
		//	//IdleConnTimeout: 30 * time.Second,
		//	//DisableCompression: true,
		//}
		//Transport: tr
		httpClient = &http.Client{}
		return httpClient
	}
	return httpClient
}

//GetRPCInfo get btcCore info
func GetRPCInfo() (string, string, string) {
	return config.BTCCore[runMode].RPCUser,
		config.BTCCore[runMode].RPCPasswd,
		config.BTCCore[runMode].getURL()

}

//GetStartBlock get start block
func GetStartBlock() int64 {
	return config.BTCCore[runMode].StartBlock
}

//GetConfirmNum get confirm num
func GetConfirmNum() int64 {
	return config.BTCCore[runMode].ConfirmNum
}
