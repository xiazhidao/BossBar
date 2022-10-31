package routers

import (
	"BossBar/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.BaseController{}, "Post:Login")
	beego.Router("/register", &controllers.BaseController{}, "Post:Register")
}
