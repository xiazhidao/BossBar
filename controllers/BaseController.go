package controllers

import (
	"BossBar/enums"
	"BossBar/models"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	curMerchant *models.Merchant //当前商户信息
}

func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{code, msg, obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}
