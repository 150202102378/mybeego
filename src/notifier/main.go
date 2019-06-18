package main

import (
	"notifier/logic"
	_ "notifier/routers"
	"time"

	"github.com/astaxie/beego"
)

func init() {

}

func notifier() {
	for {
		logic.Notifier()
		time.Sleep(time.Second * 10)
	}
}

func main() {
	logic.Notifier()
	beego.Run()
}
