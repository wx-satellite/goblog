package bootstrap

import (
	"goblog/pkg/model"
	"time"
)

func SetupDB() {
	db, err := model.ConnectDB()
	if err != nil {
		panic(err)
	}

	// GORM 底层也是使用 database/sql 来管理连接池
	sqlDB, _ := db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)

	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
