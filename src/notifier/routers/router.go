package routers

import (
	"notifier/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/update/", &controllers.UpdateController{}, "post:Update")
}
