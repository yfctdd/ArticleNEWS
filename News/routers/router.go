package routers

import (
	"News/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//beego.InsertFilter("/static/*",beego.BeforeStatic,staticFilFun)
	beego.InsertFilter("/article/*",beego.BeforeRouter,filterFun)
	beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/article/logout", &controllers.UserController{},"get:HandleLogout")
	beego.Router("/article/articleList", &controllers.ArticleController{},"get:ShowArticalList")
	beego.Router("/article/addArticle", &controllers.ArticleController{},"get:ShowAddArtical;post:HandleAddArticle")
	beego.Router("/articl/articleDetail", &controllers.ArticleController{},"get:ShowArticleDetail")
	beego.Router("/article/updateArticle", &controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	beego.Router("/article/deleteArticle", &controllers.ArticleController{},"get:DeleteArticle")
	beego.Router("/article/addArticleType", &controllers.ArticleController{},"get:ShowArticleType;post:HandleAddArticleType")
	beego.Router("/article/deleteArticleType", &controllers.ArticleController{},"get:DeleteArticleType")
}

var filterFun = func(ctx *context.Context) {
	username := ctx.Input.Session("username")

	beego.Info(",,,,,,,,,,,,,",username,",,,,,,,,,,,,,,,,,,,,")
	if username == nil {
		ctx.Redirect(302,"/login")
	}
}

//var staticFilFun = func(ctx *context.Context) {
//	ctx.Redirect(302,"/login")
//}