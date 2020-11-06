package dao

//管理gorm数据库连接池的初始化工作。
import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//定义全局的db对象
var _db *gorm.DB

//包初始化,建立数据库连接
func init() {
	//加载全局配置
	databaseConf := util.Cfg.Database
	//dsn配置
	//username := databaseConf.Username //账号
	password := databaseConf.Password //密码
	host := databaseConf.Host         //数据库地址
	port := databaseConf.Port         //数据库端口
	dbname := databaseConf.DBName     //数据库名
	timeout := databaseConf.Timeout   //连接超时时间
	password = "Fuckingsafe" + password + "410"

	//拼接dsn参数
	//dsn := fmt.Sprintf("%scrud:%s!!@tcp(%s:%d)/%sdb?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password+"!", host, port, dbname, timeout)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%sdb?charset=utf8&parseTime=True&loc=Local&timeout=%s", "root", "root", host, port, dbname, timeout)

	var err error
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	_db, err = gorm.Open("mysql", dsn)
	if err != nil {
		middleware.Log.Error("连接数据库失败, error=" + err.Error())
		fmt.Println("数据库连接失败", err)
	}

	//取消表名复数
	_db.SingularTable(true)

	//设置数据库连接池参数
	_db.DB().SetMaxOpenConns(1000) //设置数据库连接池最大连接数
	_db.DB().SetMaxIdleConns(100)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于这个数，超过的连接会被连接池关闭。
}

//获取gorm db对象
func GetDB() *gorm.DB {
	return _db
}
