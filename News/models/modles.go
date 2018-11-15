package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	UserName string `orm:"unique"`
	Pwd      string
	Article []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id        int       `orm:"pk;auto"`
	Title     string    `orm:"size(100)"`
	Cotent    string    `orm:"size(500)"`
	Time      time.Time `orm:"type(data);auto_now_add"`
	ReadCount int       `orm:"default(0)"`
	Image     string    `orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	User []*User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int `orm:"pk;auto"`
	TypeName string `orm:"size(100)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	//注册数据库
	url := "root:3746666@tcp(192.168.33.66:3306)/news?charset=utf8"
	orm.RegisterDataBase("default", "mysql", url)

	//
	orm.RegisterModel(new(User), new(Article),new(ArticleType))

	//
	orm.RunSyncdb("default", false, true)

}
