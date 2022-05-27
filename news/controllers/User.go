package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"news/models"
	"time"
)

//注册

type RegController struct {
	beego.Controller
}

func (c *RegController) ShowReg() {
	c.TplName = "register.html"
}

func (c *RegController) HandleReg() {
	//1.拿到浏览器传递的数据
	name := c.GetString("loginName")
	passw := c.GetString("password")
	//2.处理数据
	if name == "" || passw == "" {
		beego.Info("用户名密码不能为空")
		c.TplName = "register.html"
		return
	}
	//3.插入数据库
	//获取ORM对象
	o := orm.NewOrm()
	//获取插入对象
	user := models.User{}
	//插入操作
	user.UserName = name
	user.Passwd = passw

	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入数据失败")
	}
	//4.返回视图
	//c.Ctx.WriteString("注册成功")
	//重定向
	c.Redirect("/", 302)
}

//登录

type LoginController struct {
	beego.Controller
}

func (c *LoginController) ShowLogin() {
	name := c.Ctx.GetCookie("userName")

	//记住用户名
	if name != "" {
		c.Data["userName"] = name
		c.Data["check"] = "checked"
	}

	c.TplName = "login.html"
}

func (c *LoginController) HandleLogin() {

	//1.拿到浏览器传递的数据
	name := c.GetString("loginName")
	passw := c.GetString("password")
	//beego.Info(name, passw)

	//2.处理数据
	if name == "" || passw == "" {
		beego.Info("用户名密码不能为空")
		c.TplName = "login.html"
		return
	}
	//3.查找数据库
	//1.获取orm对象
	o := orm.NewOrm()
	//2.获取查询对象
	user := models.User{}
	//3.查询
	user.UserName = name
	err := o.Read(&user, "UserName")
	if err != nil {
		beego.Info("用户名错误")
		c.TplName = "login.html"
		return
	}
	//4.判断密码是否一致
	if user.Passwd != passw {
		beego.Info("密码错误")
		c.TplName = "login.html"
		return
	}
	//记住用户名的实现
	check := c.GetString("remember")
	if check == "on" {
		c.Ctx.SetCookie("userName", name, time.Second*3600)
	} else {
		c.Ctx.SetCookie("userName", "sss", -1)
	}

	c.SetSession("userName", name)
	//5.返回视图 重定向 展示用户首页主页
	c.Redirect("/Article/ShowArticle", 302)
}
