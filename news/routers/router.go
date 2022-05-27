package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"news/controllers"
)

//过滤器函数

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/")
	}
}

func init() {
	//过滤器函数                 路由规则          什么时候执行       过滤器函数
	beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFunc)
	//注册
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
	//登录
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	//展示主页
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList;post:HandleSelect")
	//添加文章
	beego.Router("/Article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	//查看详情
	beego.Router("/Content", &controllers.ArticleController{}, "get:ShowArticleContent")
	//删除文章
	beego.Router("/Article/DeleteArticle", &controllers.ArticleController{}, "get:HandleDelete")
	//编辑文章
	beego.Router("/Article/UpdateArticle", &controllers.ArticleController{}, "get:ShowUpdate;post:HandleUpdate")
	//添加分类
	beego.Router("/Article/AddArticleType", &controllers.ArticleController{}, "get:ShowAddArticleType;post:HandleArticleType")
	//用户退出
	beego.Router("/Article/Logout", &controllers.ArticleController{}, "get:Logout")
}
