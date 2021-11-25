package bootstrap

import (
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/models/user"
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"gorm.io/gorm"
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
	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)

	// 创建和维护表数据
	migration(db)
}

// gorm 默认会将键小写化作为字段名称：
//  例如，字段：Name  string `gorm:"column:name` 其实可以去掉 column

// gorm 默认是允许的 NULL，如果不要 NULL，即：
// `gorm:"not null"`

// `gorm:"-"` 设置 gorm 忽略此字段

// gorm 的自动迁移工具不支持删除数据库字段，并且只能更新字段结构，例如本来是 varchar(255) 可以修改 varchar(100)
// 如果本来有 unique ，然后在 struct 中删掉了 unique，实际表中仍旧存在，就删除无效
// 如果本来没有 index，然后在 struct 中添加了 index，那么会在表中创建索引

// 大小写不敏感，但建议使用 camelCase 风格（ 驼峰 ），比如 autoIncrement。但是写成 autoincrement 也是可以的，如果拼错则不行

// 这里有个特殊：
//	varchar(100) 改成 varchar(100);unique 仍旧不会创建唯一索引
//	只有同时将 varchar 也改了，例如 varchar(99);unique 才会创建唯一索引
func migration(db *gorm.DB) {
	err := db.AutoMigrate(new(article.Article), new(user.User), new(category.Category))
	logger.Error(err)
}
