package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type GoodsController struct {
	beego.Controller
}

func getUser(c *beego.Controller) {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Data["userName"] = ""
	} else {
		c.Data["userName"] = userName.(string)
	}
}

func (c *GoodsController) ShowIndex() {
	logs.Info("访问首页")
	getUser(&c.Controller)
	c.TplName = "index.html"
}
