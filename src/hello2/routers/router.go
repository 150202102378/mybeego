package routers

import (
	"hello2/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/check/:address", &controllers.AddrValidationController{}, "get:Get")
	beego.Router("/check", &controllers.AddrValidationController{}, "post:Post")
}
