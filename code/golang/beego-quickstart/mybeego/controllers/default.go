package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.vip111"
	c.Data["Email"] = "astaxie@gmail.com111"
	c.TplName = "index.tpl"
}
