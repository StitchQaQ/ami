package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) HelloWorld() {
	c.Ctx.WriteString("Hello World")
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("Get Method")
}

