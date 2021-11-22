package model

import (
	"fmt"
	"goblog/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB gorm的连接池对象
var DB *gorm.DB

// 初始化DB对象
func ConnectDB() (db *gorm.DB, err error) {
	gormConfig := mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
			config.GetString("database.mysql.username"),
			config.GetString("database.mysql.password"),
			config.GetString("database.mysql.host"),
			config.GetString("database.mysql.port"),
			config.GetString("database.mysql.database"),
			config.GetString("database.mysql.charset")),
	})
	// LogMode 里填写的是日志级别，分别如下：
	//   Silent ——  静默模式，不打印任何信息
	//   Error —— 发生错误了才打印
	//   Warn —— 发生警告级别以上的错误才打印
	//   Info —— 打印所有信息，包括 SQL 语句
	var level = gormlogger.Info
	if !config.GetBool("app.debug") {
		level = gormlogger.Error
	}
	// 默认使用的是 Warn
	DB, err = gorm.Open(gormConfig, &gorm.Config{
		// 日常开发，日志级别为 Warn 即可，否则命令行太多信息会反而容易让我们错过重要的信息
		Logger: gormlogger.Default.LogMode(level),
	})
	return DB, err
}
