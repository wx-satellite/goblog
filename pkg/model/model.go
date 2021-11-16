package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB gorm的连接池对象
var DB *gorm.DB

// 初始化DB对象
func ConnectDB() (db *gorm.DB, err error) {
	config := mysql.New(mysql.Config{
		DSN: "root:admin123@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})

	// LogMode 里填写的是日志级别，分别如下：
	//   Silent ——  静默模式，不打印任何信息
	//   Error —— 发生错误了才打印
	//   Warn —— 发生警告级别以上的错误才打印
	//   Info —— 打印所有信息，包括 SQL 语句
	// 默认使用的是 Warn
	DB, err = gorm.Open(config, &gorm.Config{
		// 日常开发，日志级别为 Warn 即可，否则命令行太多信息会反而容易让我们错过重要的信息
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	return DB, err
}
