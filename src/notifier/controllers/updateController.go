package controllers

import (
	"notifier/logic"

	"github.com/astaxie/beego"
)

type UpdateController struct {
	beego.Controller
}

//根据builder发过来的post请求，更新draw表
func (this *UpdateController) Update() {
	//读取post请求的内容
	body := this.Ctx.Input.RequestBody
	//处理
	logic.HandlePost(body)
}
