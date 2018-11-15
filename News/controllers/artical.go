package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"News/models"
	"math"
	"strconv"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController)Prepare()  {
	this.Data["usernameWhole"] = this.GetSession("username")
}

type ArticleController struct {
	BaseController
}



func (c *ArticleController) ShowArticalList() {
	//查询数据库 queryseter type QuerySeter interface

	//获取orm对象
	obj := orm.NewOrm()
	var articles []models.Article

	var articleTypes []models.ArticleType
	obj.QueryTable("ArticleType").All(&articleTypes)
	var queryseter orm.QuerySeter

	typeName := c.GetString("select")
	beego.Info("--------------",typeName,"------------------------")
	selStr := "全部新闻"
	if typeName == selStr ||  typeName == ""{
		c.Data["typeName"] = selStr
		queryseter = obj.QueryTable("Article")
		//	queryseter.Limit(pageSize, pageStart).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}else {
		queryseter = obj.QueryTable("Article").Filter("ArticleType__TypeName",typeName)
		c.Data["typeName"] = typeName
	}


	//queryseter := obj.QueryTable("Article")



	//queryseter.All(&articles)
	//beego.Info(articles)
	//分页
	//总记录数
	qsCount, _ := queryseter.Count()
	c.Data["qsCount"] = qsCount
	//总页数
	//单页数据数
	pageSize := 2
	pageCount := math.Ceil(float64(qsCount) / float64(pageSize))

	pageIndex, err := c.GetInt("pageIndex")

	if err != nil {
		pageIndex = 1
	}
	pageStart := pageSize * (pageIndex - 1)

	queryseter.Limit(pageSize, pageStart).RelatedSel("ArticleType").All(&articles)



	//queryseter.Limit(pageSize, pageStart).All(&articles)
	//typeName := c.GetString("select")
	//beego.Info("--------------",typeName,"------------------------")
	//selStr := "全部新闻"
	//if typeName == selStr ||  typeName == ""{
	//	c.Data["typeName"] = selStr
	////	queryseter.Limit(pageSize, pageStart).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	//}else {
	//	c.Data["typeName"] = typeName
	//}


	//


	c.Data["pageCount"] = int(pageCount)

	c.Data["articles"] = articles

	c.Data["articleTypes"] = articleTypes

	c.Data["pageIndex"] = pageIndex

	//获取查询对象
	c.Layout = "layout.html"


	c.TplName = "index.html"
}

func (c *ArticleController) ShowAddArtical() {
	var articleTypes []models.ArticleType
	obj := orm.NewOrm()
	obj.QueryTable("ArticleType").All(&articleTypes)

	c.Data["articleTypes"] = articleTypes

	c.Layout = "layout.html"

	c.TplName = "add.html"
}

func (c *ArticleController) HandleAddArticle() {
	//
	tittle := c.GetString("articleName")
	content := c.GetString("content")

	if tittle == "" || content == "" {
		c.Data["errorMeg"] = "文章标题或内容不能为空"
		c.TplName = "add.html"
		return
	}

	fp, header, err := c.GetFile("uploadname")
	if err != nil {
		c.Data["errorMeg"] = "图片上传失败"
		c.TplName = "add.html"
		return
	}
	defer fp.Close()

	//
	if header.Size > 500000 {
		c.Data["errorMeg"] = "文件太大"
		c.TplName = "add.html"
		return
	}
	//文件后缀
	pathExt := path.Ext(header.Filename)

	if pathExt != ".jpg" && pathExt != ".png" && pathExt != ".jpeg" {
		c.Data["errorMeg"] = "文件格式不正确"
		c.TplName = "add.html"
		return
	}

	//
	fileName := time.Now().Format("20060102150405") + pathExt
	//
	c.SaveToFile("uploadname", "./static/imags/"+fileName)

	//文章类型
	typeName := c.GetString("select")

	//数据库插入
	obj := orm.NewOrm()
	//获取插入对象
	var article models.Article
	var artileType models.ArticleType
	artileType.TypeName = typeName

	obj.Read(&artileType,"TypeName")

	article.Title = tittle
	article.Cotent = content
	article.Image = "/static/imags/" + fileName
	article.ArticleType = &artileType


	//插入
	id, err := obj.Insert(&article)
	if err != nil {
		c.Data["errorMeg"] = "添加文章失败"
		c.TplName = "add.html"
		return
	}

	beego.Info(id)

	c.Data["errorMeg"] = "成功"

	//c.Ctx.WriteString("成功")
	c.Redirect("/article/articleList",302)
}

func (c *ArticleController) ShowArticleDetail() {
	artId, err := c.GetInt("id")
	if err != nil {
		beego.Error("====000======")
		c.Redirect("/article/articleList", 302)
		return
	}

	obj := orm.NewOrm()
	var article models.Article
	article.Id = artId
	err = obj.Read(&article)

	var articleType models.ArticleType
	articleType.Id = article.ArticleType.Id
	obj.Read(&articleType,"Id")
	if err != nil {
		beego.Error("====000======")
		//c.Data["msg"] = "查询错误"
		c.Redirect("/article/articleList", 302)
		return
	}

	// 建立关系表数据
	m2m := obj.QueryM2M(&article,"User")
	var user models.User
	username := c.GetSession("username")
	if username == nil {
		beego.Error("====000======")
		c.Redirect("/article/articleList", 302)
		return
	}
	beego.Info("=========username========",username)
	user.UserName = username.(string)
	obj.Read(&user,"UserName")
	m2m.Add(&user)
	beego.Error("====000======")
	//查询读者与文章关系并显示
	var users []models.User
	obj.QueryTable("User").Filter("Article__Article__Id",artId).Distinct().All(&users)
	c.Data["article"] = article
	beego.Info("==========url:",article.Image)
	c.Data["articleType"] = articleType
	c.Data["users"] = users
	c.Layout = "layout.html"
	c.TplName = "content.html"

}

func (c *ArticleController) ShowUpdate() {
	artId, err := c.GetInt("id")
	if err != nil {
		c.Redirect("/article/articleList", 302)
		return
	}
	errMsg := c.GetString("errMsg")
	if errMsg != "" {
		c.Data["errMsg"] = errMsg
	}

	obj := orm.NewOrm()
	var article models.Article
	article.Id = artId
	err = obj.Read(&article)
	if err != nil {
		//c.Data["msg"] = "查询错误"
		c.Redirect("/article/articleList", 302)
		return
	}

	beego.Info(article.Image)
	c.Data["article"] = article
	c.Layout = "layout.html"
	c.TplName = "update.html"

}

func (c *ArticleController) HandleUpdate() {
	//
	id, err := c.GetInt("id")
	if err != nil {
		c.Redirect("/article/articleList", 302)
		return
	}
	tittle := c.GetString("articleName")
	content := c.GetString("content")
	if tittle == "" || content == "" {
		errMsg := "文章标题或内容不能为空"
		c.Redirect("/article/updateArticle?id="+strconv.Itoa(id)+"&errMsg="+errMsg, 302)
		return
	}

	//获取图片并更新
	fileName := updateImage(c, "uploadname")
	if fileName == "" {
		errMsg := "未上传图片"
		c.Redirect("/article/updateArticle?id="+strconv.Itoa(id)+"&errMsg="+errMsg, 302)
		return
	}

	//
	var article models.Article
	article.Id = id
	obj := orm.NewOrm()
	err = obj.Read(&article)
	if err != nil {
		errMsg := "原文章不存在"
		c.Redirect("/article/updateArticle?id="+strconv.Itoa(id)+"&errMsg="+errMsg, 302)
		return
	}
	article.Title = tittle
	article.Cotent = content
	article.Image = fileName
	//article.Time = time.Now()

	iNum, err := obj.Update(&article)
	if err != nil {
		beego.Error("更新失败", err)
		errMsg := "更新失败"
		c.Redirect("/article/updateArticle?id="+strconv.Itoa(id)+"&errMsg="+errMsg, 302)
		return

	}
	beego.Info("iNum", iNum)

	c.Redirect("/article/articleList", 302)

}

func updateImage(c *ArticleController, url string) string {
	fp, header, err := c.GetFile("uploadname")
	if err != nil {
		c.Data["errorMeg"] = "图片上传失败"
		c.TplName = "add.html"
		return ""
	}
	defer fp.Close()

	//
	if header.Size > 500000 {
		c.Data["errorMeg"] = "文件太大"
		c.TplName = "add.html"
		return ""
	}
	//文件后缀
	pathExt := path.Ext(header.Filename)

	if pathExt != ".jpg" && pathExt != ".png" && pathExt != ".jpeg" {
		c.Data["errorMeg"] = "文件格式不正确"
		c.TplName = "add.html"
		return ""
	}

	//
	fileName := time.Now().Format("2006-01-02 15:04:05") + pathExt
	//
	c.SaveToFile("uploadname", "./static/imags/"+fileName)
	return "./static/imags/" + fileName
}

func (c *ArticleController) DeleteArticle() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Error(err)
		c.Data["msg"] = "请求路径错误"
		c.Redirect("/article/articleList", 302)
		return
	}

	var article models.Article
	article.Id = id
	obj := orm.NewOrm()
	_, err = obj.Delete(&article)
	if err != nil {
		beego.Error(err)
		c.Data["msg"] = "删除错误"
		c.Redirect("/article/articleList", 302)
		return
	}

	c.Redirect("/article/articleList", 302)

}

func (c *ArticleController) ShowArticleType() {
	errMsg := c.GetString("errMsg")
	if errMsg != "" {
		c.Data["errMsg"] = errMsg
	}

	var articleTypes []models.ArticleType
	obj := orm.NewOrm()
	_, err := obj.QueryTable("ArticleType").All(&articleTypes)
	if err != nil {
		c.Data["errMsg"] = "获取信息失败"
		c.TplName = "addType.html"
		return
	}
	c.Data["articleTypes"] = articleTypes
	c.Layout = "layout.html"
	c.TplName = "addType.html"

}

func (c *ArticleController) HandleAddArticleType() {
	//获取类别名称
	typeName := c.GetString("typeName")
	if typeName == "" {
		beego.Error("获取类型名错误")
		errMsg := "获取类型名错误"
		c.Redirect("/article/addArticleType?errMsg="+errMsg, 302)
		return
	}
	//存入数据库
	var articleType models.ArticleType
	articleType.TypeName = typeName
	obj := orm.NewOrm()
	_, err := obj.Insert(&articleType)
	if err != nil {
		beego.Error("新建类型名错误")
		errMsg := "新建类型名错误"
		c.Redirect("/article/addArticleType?errMsg="+errMsg, 302)
		return
	}
	//重定向返回
	c.Redirect("/article/addArticleType", 302)

}

func (c *ArticleController) DeleteArticleType() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Error("无法获取删除信息")
		errMsg := "无法获取删除信息"
		c.Redirect("/article/addArticleType?errMsg="+errMsg, 302)
		return
	}

	var articleType models.ArticleType
	articleType.Id = id
	obj := orm.NewOrm()
	_, err = obj.Delete(&articleType)
	if err != nil {
		beego.Error("无法删除信息")
		errMsg := "无法删除信息"
		c.Redirect("/article/addArticleType?errMsg="+errMsg, 302)
		return
	}

	c.Redirect("/article/addArticleType", 302)
}
