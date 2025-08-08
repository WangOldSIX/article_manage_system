package routers

import (
	"github.com/beego/beego/v2/server/web/context"
	"day_day_fresh/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeRouter, filterFunc)
	beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	//激活用户
	beego.Router("/active", &controllers.UserController{},"get:ActiveUser")
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")

	//跳转首页
	beego.Router("/", &controllers.UserController{},"get:ShowLogin")

	//退出登录
	beego.Router("/user/logout", &controllers.UserController{},"get:Logout")

	//用户中心
	beego.Router("/user/userCenterInfo", &controllers.UserController{},"get:ShowUserCenterInfo")
	//全部订单
	beego.Router("/user/userCenterOrder", &controllers.UserController{},"get:ShowUserCenterOrder")
	//地址页面
	beego.Router("/user/userCenterSite", &controllers.UserController{},"get:ShowUserCenterSite;post:HandleUserCenterSite")

}

var filterFunc=func(ctx *context.Context) {
	// 过滤器
	// 如：判断是否登录，是否有权限等
	// 若用户未登录，则重定向到登录页面
	// 若用户无权限，则重定向到无权限页面
	// 若用户登录，则放行
	userName:=ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302,"/login")
		return
	}
}