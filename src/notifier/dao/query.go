package dao

import (
	"notifier/models"

	"github.com/astaxie/beego/orm"
)

func RechargeQuery() []*models.ConfirmedTx {
	o := orm.NewOrm()

	//查询status为0的数据并保存在confirms
	var confirms []*models.ConfirmedTx
	qs := o.QueryTable("confirmed_tx")
	//qs := o.QueryTable("confirmedTx")
	qs.Filter("Status", 0).All(&confirms)
	return confirms
}

func WithdrawQuery() []*models.Draw {
	o := orm.NewOrm()

	//查询status为0的数据并保存在draws
	var draws []*models.Draw
	qs := o.QueryTable("draw")
	qs.Filter("Status", 0).All(&draws)
	return draws
}
