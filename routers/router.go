package routers

import (
	"day_day_fresh/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	//激活用户
	beego.Router("/active", &controllers.UserController{},"get:ActiveUser")
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	
}
