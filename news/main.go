package main

import (
	"github.com/astaxie/beego"
	_ "news/models" //_为启动时调用包里的init函数
	_ "news/routers"
	"strconv"
)

func main() {
	//视图函数映射要在run之前完成
	beego.AddFuncMap("ShowPrePage", HandlePrePage)
	beego.AddFuncMap("ShowNextPage", HandleNextPage)

	beego.Run()

}

//视图函数

//上一页

func HandlePrePage(data int) string {

	pageIndex := data - 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}

//下一页

func HandleNextPage(data int) string {

	pageIndex := data + 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
