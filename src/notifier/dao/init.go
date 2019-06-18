package dao

import (
	"notifier/models"
	"github.com/astaxie/beego/orm"
)

func init() {
	//注册驱动
	orm.RegisterDriver("postgres", orm.DRPostgres)
	//注册数据库
	orm.RegisterDataBase("default", "postgres", "user=te password=traveltao dbname=te host=127.0.0.1 port=5432 sslmode=disable")
	//orm.RegisterDataBase("default", "postgres", "user=te password=traveltao dbname=te host=192.168.1.57 port=55555 sslmode=disable")
	//注册模型
	orm.RegisterModel(new(models.ConfirmedTx), new(models.Draw))
}
