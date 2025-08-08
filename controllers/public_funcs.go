package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

func getUser(c *beego.Controller) string{
	userName := c.GetSession("userName")
	if userName == nil {
		c.Data["userName"] = ""
	} else {
		c.Data["userName"] = userName.(string)
	}
	//不加逻辑判断，因为如果不登陆访问不到这里
	return userName.(string)
}
