package controllers

import (
	"github.com/astaxie/beego"
	"News/models"
	"github.com/astaxie/beego/orm"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

//展示注册页面
func (c *UserController)ShowRegister()  {
	c.TplName = "register.html"
}

//处理注册数据
func (c *UserController)HandleRegister()  {
	username := c.GetString("userName")
	pwd := c.GetString("password")
	//
	if username == "" || pwd == "" {
		c.Data["errorMsg"] = "用户名或密码不能为空"
		c.TplName = "register.html"
		return
	}

	var user models.User
	obj := orm.NewOrm()

	user.UserName = username
	user.Pwd = pwd


	id,err := obj.Insert(&user)
	if err != nil {
		c.Data["errorMsg"] = "注册失败,请重新输入"
		c.TplName = "register.html"
		return
	}
	beego.Info(id)


	//c.Ctx.WriteString("注册成功")
	c.Redirect("/login",302)

}

func (c *UserController)ShowLogin()  {
	etc := c.Ctx.GetCookie("username")
	etcByte,_ := base64.StdEncoding.DecodeString(etc)
	username := string(etcByte)
	if username != "" {
		c.Data["username"] = username
		c.Data["checked"] = "checked"
	}else {
		c.Data["username"] = ""
		c.Data["checked"] = ""
	}

	//c.SetSession("username",username)
	c.TplName = "login.html"
}

func (c *UserController)HandleLogin()   {
	username := c.GetString("userName")
	pwd := c.GetString("password")

	if username == "" || pwd == "" {
		c.Data["errorMsg"] = "用户名或密码不能为空"
		c.TplName = "login.html"
		return
	}

	remember := c.GetString("remember")
	//beego.Info("==========",remember)
	etc := base64.StdEncoding.EncodeToString([]byte(username))
	if remember == "on" {
		c.Ctx.SetCookie("username",etc,3600 * 1)
	}else {
		c.Ctx.SetCookie("username",etc, -1)
	}

	var user models.User
	obj := orm.NewOrm()

	user.UserName = username

	err := obj.Read(&user,"UserName")

	if err != nil {
		c.Data["errorMsg"] = "用户名不存在，请重新输入"
		c.TplName = "register.html"
		return
	}

	if user.Pwd != pwd {
		c.Data["errorMsg"] = "密码错误"
		c.TplName = "register.html"
		return
	}

	c.SetSession("username",username)
	//c.Ctx.WriteString("登陆成功")

	c.Redirect("/article/articleList",302)
}

func (c *UserController)HandleLogout()  {
	//
	c.DelSession("username")
	c.Redirect("/login",302)
}