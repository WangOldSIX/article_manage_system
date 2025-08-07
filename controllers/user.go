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
	userPwd := c.Ctx.GetCookie("userpwd")
	logs.Info("加密的用户名,密码:\n%v\n%v", userName, userPwd)
	//解码
	temp, _ := base64.StdEncoding.DecodeString(userName)
	temp2, _ := base64.StdEncoding.DecodeString(userPwd)
	if string(temp) == "" {
		c.Data["userName"] = ""
		c.Data["checked"] = ""
	} else {
		c.Data["userName"] = string(temp)
		c.Data["checked"] = "checked"
	}
	if string(temp2) == "" {
		c.Data["userPwd"] = ""
		c.Data["checked"] = ""
	} else {
		c.Data["userPwd"] = string(temp2)
		c.Data["checked"] = "checked"
	}

	if userName != "" {
		c.Data["remember_username"] = "checked"
	}
	if userPwd != "" {
		c.Data["remember_pwd"] = "checked"
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
	remember_pwd := c.GetString("remember_pwd")
	if remember_username == "on" {
		tmp_username := base64.StdEncoding.EncodeToString([]byte(remember_username))
		logs.Info("记住用户名(加密):%v", tmp_username)
		c.Ctx.SetCookie("username", tmp_username, 3600*24*7)
	} else {
		c.Ctx.SetCookie("username", "", -1)
	}
	if remember_pwd == "on" {
		tmp_pwd := base64.StdEncoding.EncodeToString([]byte(remember_username))
		logs.Info("记住密码(加密):%v", tmp_pwd)
		c.Ctx.SetCookie("userpwd", tmp_pwd, 3600*24*7)
	} else {
		c.Ctx.SetCookie("userpwd", "", -1)
	}
	c.Data["userName"] = user.Name
	c.Data["userPwd"] = pwd
	c.Data["errmsg"] = "登录成功"

	//配置session
	c.SetSession("userName", userName)
	c.SetSession("userPwd", user.PassWord)

	c.Redirect("/", 302)
}

// 退出登录
func (c *UserController) Logout() {
	c.DelSession("userName")
	c.DelSession("userPwd")
	//跳转视图
	// c.Redirect("/login", 302)
	c.Redirect("/", 302)
}

func (c *UserController) ShowUserCenterInfo() {
	getUser(&c.Controller)
	//获取用户
	userName := c.GetSession("userName")
	c.Data["userName"] = userName
	//查询用户信息
	o := orm.NewOrm()
	//高级查询，表关联
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName).Filter("Isdefault", true).One(&addr)
	if addr.Id == 0 {
		c.Data["addr"] = ""
	} else {
		c.Data["addr"] = addr
	}
	c.Layout = "userCenterLayout.html"
	c.TplName = "user_center_info.html"
}

// 用户中心订单
func (c *UserController) ShowUserCenterOrder() {
	getUser(&c.Controller)
	c.Layout = "userCenterLayout.html"
	c.TplName = "user_center_order.html"
}

// 用户中心收货地址
func (c *UserController) ShowUserCenterSite() {
	userName := getUser(&c.Controller)
	// c.Data["userName"] = userName
	//查询用户地址
	o := orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").
	RelatedSel("User").
	Filter("User__Name", userName).
	Filter("Isdefault", true).
	One(&addr)

	logs.Info("用户:%v 默认地址:%v", userName, addr)
	c.Data["addr"] = addr

	c.Layout = "userCenterLayout.html"
	c.TplName = "user_center_site.html"
}

func (c *UserController) HandleUserCenterSite() {
	//获取数据
	receiver := c.GetString("receiver")
	addr := c.GetString("addr")
	zipcode := c.GetString("zipCode")
	phone := c.GetString("phone")
	logs.Info("收件人:%v 地址:%v 邮编:%v 手机:%v", receiver, addr, zipcode, phone)
	//检验数据
	if receiver==""{
		c.Data["errmsg"] = "收件人不能为空"
		logs.Error("收件人不能为空")
		c.Redirect("/user/userCenterSite", 302)
		return
	}
	if addr==""{
		c.Data["errmsg"] = "地址不能为空"
		logs.Error("地址不能为空")
		c.Redirect("/user/userCenterSite", 302)
		return
	}
	if zipcode==""{
		c.Data["errmsg"] = "邮编不能为空"
		logs.Error("邮编不能为空")
		c.Redirect("/user/userCenterSite", 302)
		return
	}
	if phone==""{
		c.Data["errmsg"] = "手机不能为空"
		logs.Error("手机不能为空")
		c.Redirect("/user/userCenterSite", 302)
		return
	}
	//处理数据
	//插入操作
	o := orm.NewOrm()
	var addrUser models.Address
	addrUser.Isdefault = true
	err := o.Read(&addrUser, "Isdefault")
	//添加默认地址之前，需要把原来的默认地址取消默认
	if err == nil {
		addrUser.Isdefault = false
		o.Update(&addrUser)
	}
	//更新默认地址，给原来的地址对象赋值了，这时候用原来的地址对象再插入就不对了
	//关联对象
	userName := c.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	var addrUserNew models.Address
	addrUserNew.Receiver = receiver
	addrUserNew.Addr = addr
	addrUserNew.Zipcode = zipcode
	addrUserNew.Phone = phone
	addrUserNew.Isdefault = true
	//关联的用户
	addrUserNew.User = &user
	o.Insert(&addrUserNew)

	logs.Info("用户:%v 添加收货地址%v成功", user, addrUserNew)
	//返回视图
	c.Redirect("/user/userCenterSite", 302)
}
