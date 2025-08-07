package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type GoodsController struct {
	beego.Controller
}

func (c *GoodsController) ShowIndex() {
	logs.Info("访问首页")
	userName:=c.GetSession("userName")
	if userName==nil{
		c.Data["userName"]=""
	}else{
		c.Data["userName"]=userName.(string)
	}

	c.TplName = "index.html"
}


