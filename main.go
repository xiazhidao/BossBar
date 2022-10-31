package main

import (
	_ "BossBar/routers"
	_ "BossBar/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.AppName = beego.AppConfig.String("frontend::appname")
	beego.BConfig.Listen.HTTPPort, _ = beego.AppConfig.Int("frontend::httpport")
	beego.BConfig.WebConfig.Session.SessionOn = false
	beego.Run()
}
