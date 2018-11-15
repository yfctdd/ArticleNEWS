package main

import (
	_ "News/routers"
	"github.com/astaxie/beego"
	_ "News/models"
)

func main() {
	beego.AddFuncMap("PrePage",PrePage)
	beego.AddFuncMap("NextPage",NextPage)
	beego.Run()
}

func PrePage(pageIndex int)int{
	if pageIndex >= 2 {
		return pageIndex - 1
	}
	return 1
}

func NextPage(pageIndex,pageCount int) int {
	if pageIndex + 1 > pageCount {
		return pageIndex
	}
	return pageIndex + 1
}

