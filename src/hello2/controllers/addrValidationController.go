package controllers

import (
	"hello2/addrvalidation"

	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
)

type AddrValidationController struct {
	beego.Controller
}

func (this *AddrValidationController) Prepare() {

}

func (this *AddrValidationController) Post() {

	//获取地址
	reqBody := this.Ctx.Input.RequestBody
	address := gjson.GetBytes(reqBody, "address").String()

	//判断地址
	addrvali := addrvalidation.GetAddrValidation()
	result := addrvali.Verify(address, "BitcoinPrefix")
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *AddrValidationController) Get() {
	address := this.GetString(":address", "")
	addrvali := addrvalidation.GetAddrValidation()
	result := addrvali.Verify(address, "BitcoinPrefix")
	this.Data["json"] = &result
	this.ServeJSON()
}
