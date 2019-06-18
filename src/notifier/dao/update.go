package dao

import (
	"fmt"
	"notifier/models"

	"github.com/astaxie/beego/orm"
)

func RechargeUpdate(confirm *models.ConfirmedTx, status int) bool {
	o := orm.NewOrm()
	//fmt.Println(status)
	switch status {
	case 1:
		confirm.Status = 1
	case 0:
		confirm.Status = 2
	}
	_, err := o.Update(confirm, "Status")
	if err == nil {
		fmt.Println("更新成功")
		return true
	} else {
		fmt.Println("更新失败")
		return false
	}
}

func WithdrawUpdate(draws *models.Draw, status int) bool {
	o := orm.NewOrm()
	switch status {
	case 2:
		draws.Status = 2
	case 1:
		draws.Status = 1
	//尚未确认失败后取何值
	case 0:
		draws.Status = 3
	}
	_, err := o.Update(draws, "Status")
	if err == nil {
		fmt.Println("更新成功")
		return true
	} else {
		fmt.Println("更新失败")
		return false
	}
}
