package controllers

import (
	constvar "day_day_fresh/const_var"
	"day_day_fresh/models"
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) ShowRegister() {
	c.TplName = "register.html"
}

/**
 * @Title HandleRegister
 * @Description 处理注册逻辑
 */
func (c *UserController) HandleRegister() {
	//TODO: handle register logic here
	//获取数据
	userName := c.GetString("user_name")
	pwd := c.GetString("pwd")
	cpwd := c.GetString("cpwd")
	email := c.GetString("email")
	//校验数据
	if len(userName) == 0 || len(pwd) == 0 || len(cpwd) == 0 || len(email) == 0 {
		c.Data["errmsg"] = "用户名、密码、确认密码、邮箱不能为空"
		logs.Error("用户名、密码、确认密码、邮箱不能为空")
		c.TplName = "register.html"
		return
	}
	if pwd != cpwd {
		c.Data["errmsg"] = "两次密码不一致"
		logs.Error("两次密码不一致")
		c.TplName = "register.html"
		return
	}
	//校验邮箱格式
	if !regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$").MatchString(email) {
		c.Data["errmsg"] = "邮箱格式不正确"
		logs.Error("邮箱格式不正确")
		c.TplName = "register.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.PassWord = pwd
	user.Email = email
	logs.Info("%v %v %v ", user.Name, user.PassWord, user.Email)
	_, err := o.Insert(&user)
	if err != nil {
		c.Data["errmsg"] = "注册失败，用户名已存在"
		logs.Error("注册失败，用户名已存在:%v", err)
		c.TplName = "register.html"
		return
	}
	logs.Info("注册成功，用户:%v", user)
	//发送邮件，激活用户
	emailConfig := constvar.EmailCofig
	logs.Info("邮箱配置:%v", emailConfig)
	emailConn := utils.NewEMail(emailConfig)
	if emailConn == nil {
		logs.Error("邮箱连接失败")
		c.Data["errmsg"] = "注册失败，邮箱连接失败"
		c.TplName = "register.html"
		return
	}
	emailConn.From = "1447895999@qq.com"
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜系统注册激活链接"
	emailConn.Text = constvar.Ipv4_local_ASUS + ":80/active?id=" + strconv.Itoa(user.Id) //发给用户的邮件内容：激活连接
	logs.Info("邮件内容%v", emailConn.Text)
	err = emailConn.Send()
	if err != nil && !strings.Contains(err.Error(), "short response") {
		logs.Error("发送邮件失败:%v", err)
		c.Data["errmsg"] = "注册失败，发送邮件失败"
		c.TplName = "register.html"
		return
	}
	logs.Info("发送邮件成功")

	//返回视图
	c.Ctx.WriteString("注册成功，请前往邮箱激活账号")
}

// 激活处理
func (c *UserController) ActiveUser() {
	//获取用户id
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("获取用户id失败:%v", err)
		c.Data["errmsg"] = "激活失败，获取用户id失败"
		c.TplName = "register.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	err = o.Read(&user)
	if err != nil {
		logs.Error("读取用户信息失败:%v", err)
		c.Data["errmsg"] = "激活失败，读取用户信息失败"
		c.TplName = "register.html"
		return
	}
	if user.Active {
		logs.Info("用户:%v 已经激活", user)
		c.Data["errmsg"] = "用户已经激活"
		c.TplName = "register.html"
		return
	}
	user.Active = true
	_, err = o.Update(&user)
	if err != nil {
		logs.Error("更新用户激活状态失败:%v", err)
		c.Data["errmsg"] = "激活失败，更新用户激活状态失败"
		c.TplName = "register.html"
		return
	}
	logs.Info("用户:%v 激活成功", user)
	c.Data["errmsg"] = "激活成功"
	c.TplName = "register.html"
	// c.Redirect("/login", 302)
}

func (c *UserController) ShowLogin() {
	userName := c.Ctx.GetCookie("username")
	logs.Info("加密的用户名\n%v",userName)
	//解码
	temp,_:=base64.StdEncoding.DecodeString(userName)
	if string(temp) == "" {
		c.Data["userName"]=""
		c.Data["checked"]=""
	}else{
		c.Data["userName"]=string(temp)
		c.Data["checked"]="checked"
	}

	if userName != "" {
		c.Data["remember_username"] = "checked"
	}
	c.TplName = "login.html"
}



func (c *UserController) HandleLogin() {
	userName := c.GetString("username")
	pwd := c.GetString("pwd")
	if len(userName) == 0 || len(pwd) == 0 {
		logs.Error("用户名或密码不能为空")
		c.Data["errmsg"] = "用户名或密码不能为空"
		c.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		logs.Error("读取用户信息失败:%v", err)
		c.Data["errmsg"] = "登录失败，用户名不存在"
		c.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		logs.Info("用户名或密码错误:%v", user)
		c.Data["errmsg"] = "用户名或密码错误"
		c.TplName = "login.html"
		return
	}
	if !user.Active {
		logs.Info("用户:%v 未激活", user)
		c.Data["errmsg"] = "用户未激活,请前往邮箱激活账号"
		c.TplName = "login.html"
		return
	}

	logs.Info("用户:%v 登录成功", user)
	//记住用户选项
	remember_username := c.GetString("remember_username")
	if remember_username == "on" {
		tmp_username:= base64.StdEncoding.EncodeToString([]byte(remember_username))
		logs.Info("记住用户名(加密):%v", tmp_username)
		c.Ctx.SetCookie("username", tmp_username, 3600*24*7)
	} else {
		c.Ctx.SetCookie("username", "", -1)
	}
	c.Data["userName"] = user.Name
	c.Data["errmsg"] = "登录成功"

	//配置session
	c.SetSession("userName",userName)

	c.Redirect("/", 302)
}
