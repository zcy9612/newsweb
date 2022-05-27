package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//表的设计

//用户表
//用户与文章是多对多关系 多对多orm创建一个表用于保存两个表之间的关系

type User struct {
	Id       int //默认为主键
	UserName string
	Passwd   string
	Articles []*Article `orm:"rel(m2m)"` //多对多的表创立外键
}

//文章表
//文章表与文章类型表是n ：1关系

type Article struct {
	Id2     int          `orm:"pk;auto"`                     //pk设置主键；auto编号自动增加
	Title   string       `orm:"size(20)"`                    //标题 size(20)设置字段大小
	Content string       `orm:"size(500)"`                   //内容
	Img     string       `orm:"size(50);null"`               //图片（路径） null允许为空，默认为非空
	Time    time.Time    `orm:"type(datetime);auto_now_add"` //发布时间 type(datetime)时间，第一次保存时设置时间
	Count   int          `orm:"default(0)"`                  //阅读量 默认为0
	AType   *ArticleType `orm:"rel(fk)"`                     //设置表外键
	Users   []*User      `orm:"reverse(many)"`
}

//多表关系
//文章类型表

type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Article  []*Article `orm:"reverse(many)"` //设置表外键与`orm:"rel(fk)""`成对出现
}

//ORM映射数据库初始化
func init() {
	//1.连接数据库参数： 数据库别名 驱动名称 连接字符串"用户名:密码@传输协议(主机ip:主机端口号)/数据库名称?编码格式"
	orm.RegisterDataBase("default",
		"mysql",
		"root:root@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	//2.注册表 可以new多个    new(表名)
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	//3.生成表         参数数据库别名 是否强制更新   是否能看见创建表的过程
	orm.RunSyncdb("default", false, true)
}
