package controllers

import(
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *UserController) HandleRegister() {
	//TODO: handle register logic here
}
