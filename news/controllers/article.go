package controllers

import (
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"news/models"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) HandleSelect() {

}

//文章列表页

func (c *ArticleController) ShowArticleList() {
	//判断用户登录
	userName := c.GetSession("userName")

	//显示文章列表
	//1.高级查询
	o := orm.NewOrm()
	//查询Article表
	qs := o.QueryTable("Article")
	var articles []models.Article
	qs.All(&articles) //select * from Article 查询出来的内容放入articles对象中

	pageIndex := c.GetString("pageIndex")
	typeName := c.GetString("select")

	pageIndex1, err := strconv.Atoi(pageIndex)
	if err != nil {
		pageIndex1 = 1 //默认设置为1
	}
	var count int64
	//获取总条目
	if typeName == "" {
		//默认访问时给出所有类型的条目
		count, err = qs.RelatedSel("AType").Count()
		if err != nil {
			beego.Info("查询错误")
			return
		}
	} else {
		//查询表中对应的类型有多少条记录    article表中的字段        对应的表中的字段  根据查询的对象
		count, err = qs.RelatedSel("AType").Filter("AType__TypeName", typeName).Count()
	}
	//获取总页数
	pageSize := 2 //页面显示条数

	//分页
	start := pageSize * (pageIndex1 - 1) //起始位置
	//参数：一页限制多少条记录，起始位置
	qs.Limit(pageSize, start).RelatedSel("AType").All(&articles)

	pageCount := float64(count) / float64(pageSize)
	pageCount1 := math.Ceil(pageCount)

	//首页末页处理
	firstPage := false
	endPage := false
	if pageIndex1 == 1 {
		firstPage = true
	}
	if pageIndex1 == int(pageCount1) {
		endPage = true
	}

	//获取新闻类型
	var types []models.ArticleType //定义一个存储所有types的对象宿主
	conn, err := redis.Dial("tcp", ":6379")
	//types是我们自定义的结构体 不能使用redis.values scan赋值
	//需要借助序列化和反序列化进行处理
	buffer, err := redis.Bytes(conn.Do("get", "types"))
	if err != nil {
		beego.Info("获取redis数据失败")
		return
	}

	dec := gob.NewDecoder(bytes.NewReader(buffer))
	dec.Decode(&types)
	beego.Info(types)

	if len(types) == 0 {
		//如果redis中没有数据从常规的数据库中读取数据
		o.QueryTable("ArticleType").All(&types)
		//将数据存入redis中
		//使用redis存新闻类型
		//序列化
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(types)

		//将序列化后的数据存入redis 注：传入的类型为bytes类型
		_, err = conn.Do("set", "types", buffer.Bytes())
		if err != nil {
			beego.Info("redis数据库操作错误，存入数据失败")
			return
		}
		beego.Info("从mysql中获取数据")
	}
	c.Data["types"] = types

	//根据类型获取数据

	var articleswithTYpe []models.Article

	//beego.Info(typeName)
	//处理数据
	if typeName == "" {
		//beego.Info("下拉框传递数据失败")
		qs.Limit(pageSize, start).RelatedSel("AType").All(&articleswithTYpe)
	} else {
		qs.Limit(pageSize, start).RelatedSel("AType").Filter("AType__TypeName", typeName).All(&articleswithTYpe)
	}

	//beego.Info(count)
	c.Data["typeName"] = typeName
	c.Data["count"] = count          //总条目
	c.Data["pageCount"] = pageCount1 //总页数
	c.Data["pageIndex"] = pageIndex1 //当前页
	c.Data["FirstPage"] = firstPage
	c.Data["EndPage"] = endPage
	c.Data["userName"] = userName
	//beego.Info(articles[0])
	//2.把数据传递给视图显示
	//传数据
	c.Data["articles"] = articleswithTYpe
	//使用layout设置页面布局
	c.Layout = "layout.html"
	c.TplName = "index.html"
}

func (c *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()

	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	c.Data["types"] = types

	c.Layout = "layout.html"
	c.TplName = "add.html"
}

/*
	1.拿到浏览器传递的数据
	2.处理数据
	3.查找数据库
	    获取orm对象
	    获取查询对象
	    查询
	    判断密码是否一致
	5.返回视图
*/
func (c *ArticleController) HandleAddArticle() {
	//1.拿到浏览器传递的数据
	Title := c.GetString("title")     //拿标题
	Content := c.GetString("content") //拿内容
	f, h, err := c.GetFile("upload")  //拿图片
	if err != nil {
		beego.Info("上传文件失败")
		return
	}
	defer f.Close()

	//2.处理数据

	//1.判断文件格式
	ext := path.Ext(h.Filename)
	beego.Info(ext)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("上传文件格式不正确")
		return
	}
	//2.判断文件大小
	if h.Size > 5000000 {
		beego.Info("上传文件太大")
		return
	}
	//3.不能重名
	//文件名不能使用英文:当作文件名称
	filename := time.Now().Format("2006-01-02 15：04：05")

	//将文件保存至本地服务器路径下 视图input中的name   路径
	err = c.SaveToFile("upload", "./static/images/"+filename+ext)
	if err != nil {
		beego.Info("下载失败", err)
		return
	}

	//beego.Info(filename + ext)
	//beego.Info(Title, Content)
	//	3.查找数据库
	//	    获取orm对象
	//	    创建插入对象
	//	    插入
	//	    判断密码是否一致
	o := orm.NewOrm()
	article := models.Article{}
	article.Title = Title
	article.Content = Content
	article.Img = "./static/images/" + filename + ext

	//类型数据插入
	typeName := c.GetString("select")
	beego.Info(typeName)
	if typeName == "" {
		beego.Info("下拉框数据获取失败")
		return
	}
	//创建一个articleTye对象
	var articleTye models.ArticleType
	//将页面中拿到的数据进行赋值
	articleTye.TypeName = typeName

	err = o.Read(&articleTye, "TypeName")
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	article.AType = &articleTye

	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("插入失败", err)
		return
	}

	//5.返回视图
	c.Redirect("/Article/ShowArticle", 302)
}

func (c *ArticleController) ShowArticleContent() {
	//接收来自前端页面的id值
	id := c.GetString("id")
	//beego.Info(id)
	//根据id查询
	o := orm.NewOrm()
	id2, _ := strconv.Atoi(id)
	article := models.Article{Id2: id2}
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询失败")
		return
	}
	//点击查看详情 计数加一
	article.Count += 1
	o.Update(&article)
	//多表查询新闻类型
	typeID := article.AType
	var articleType models.ArticleType
	o.QueryTable("ArticleType").Filter("id", typeID).One(&articleType)
	c.Data["Atype"] = articleType

	//多对多插入读者
	//在Article表中插入Users字段
	//获取操作对象
	artile := models.Article{Id2: id2}
	//获取多对多操作对象
	m2m := o.QueryM2M(&artile, "Users")
	//获取插入对象
	userName := c.GetSession("userName")
	//beego.Info(userName)
	user := models.User{}
	user.UserName = userName.(string)
	o.Read(&user, "UserName")
	//多对多插入
	_, err = m2m.Add(&user)
	if err != nil {
		beego.Info("多对多插入失败")
		return
	}
	o.Update(&article)

	//多对多查询
	//1.LoadRelated查询会有重复的数据
	//o.LoadRelated(&article, "Users")
	//2.第二种采用过滤器多对多查询          表名  字段名    表名      查询内容的字段       去重
	//o.QueryTable("Article").Filter("Users__User__UserName", userName.(string)).Distinct().Filter("Id2", id2).One(&article)
	//多对多查询
	//从User表中查询
	var users []models.User
	// 采用过滤器多对多查询             表名           表名    字段名  表名  查询内容的字段  去重
	o.QueryTable("User").Filter("Articles__Article__Id2", id2).Distinct().All(&users)
	//beego.Info(article)

	c.Data["article"] = article
	c.Data["users"] = users

	c.Layout = "layout.html"
	c.TplName = "content.html"
}

//url传值
//delete

func (c *ArticleController) HandleDelete() {
	id, _ := c.GetInt("id")

	o := orm.NewOrm()
	article := models.Article{Id2: id}
	o.Delete(&article)

	c.Redirect("/Article/ShowArticle", 302)
}

func (c *ArticleController) ShowUpdate() {
	//接收来自前端页面的id值
	id := c.GetString("id")
	//beego.Info(id)
	if id == "" {
		beego.Info("连接错误")
		return
	}
	//根据id查询
	o := orm.NewOrm()
	id2, _ := strconv.Atoi(id)
	article := models.Article{Id2: id2}
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询失败")
		return
	}

	c.Data["article"] = article

	c.Layout = "layout.html"
	c.TplName = "update.html"
}

func (c *ArticleController) HandleUpdate() {
	//拿数据
	title := c.GetString("title")
	content := c.GetString("content")
	//判断
	if title == "" || content == "" {
		beego.Info("更新数据失败")
		return
	}

	//拿文件
	f, h, err := c.GetFile("upload")
	if err != nil {
		beego.Info("上传文件失败")
		return
	}
	defer f.Close()

	if h.Size > 5000000 {
		beego.Info("图片过大")
		return
	}

	ext := path.Ext(h.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		beego.Info("文件格式错误")
		return
	}
	filename := time.Now().Format("2006-01-02 15：04：05")
	err = c.SaveToFile("upload", "./static/images/"+filename+ext)
	if err != nil {
		beego.Info("下载失败", err)
		return
	}
	//更新 先根据id将该信息取出来，在进行赋值，最后更新数据库
	id, _ := c.GetInt("id")
	o := orm.NewOrm()
	article := models.Article{Id2: id}
	err = o.Read(&article)
	if err != nil {
		beego.Info("文章读取失败")
		return
	}

	article.Content = content
	article.Title = title
	article.Img = "./static/images/" + filename + ext

	_, err = o.Update(&article)
	if err != nil {
		beego.Info("更新失败")
		return
	}
	c.Redirect("/Article/ShowArticle", 302)
}
func (c *ArticleController) ShowAddArticleType() {
	//读取类型表显示
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	qs := o.QueryTable("ArticleType")
	_, err := qs.All(&articleTypes)
	if err != nil {
		beego.Info("读取ArticleType表失败")
	}
	c.Data["types"] = articleTypes

	c.Layout = "layout.html"
	c.TplName = "addType.html"
}
func (c *ArticleController) HandleArticleType() {
	typeName := c.GetString("type")
	beego.Info(typeName)
	if typeName == "" {
		beego.Info("添加数据为空")
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	_, err := o.Insert(&articleType)
	if err != nil {
		beego.Info("插入文章类型失败")
		return
	}
	c.Redirect("/Article/AddArticleType", 302)
}

//退出登录

func (c *ArticleController) Logout() {
	//删除session
	c.DelSession("userName")
	//跳转到登录
	c.Redirect("/", 302)
}
