package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	if err := initDB(); err != nil {
		panic(err)
	}
	if err := createTables(); err != nil {
		panic(err)
	}

}

func initDB() (err error) {
	config := mysql.Config{
		User:                 "root",
		Passwd:               "admin123",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}
	DB, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return
	}
	// 设置最大空闲连接数
	DB.SetMaxIdleConns(25)
	// 设置最大连接数
	DB.SetMaxOpenConns(25)
	DB.SetConnMaxLifetime(5 * time.Hour)
	return DB.Ping()
}

func createTables() (err error) {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `
	_, err = DB.Exec(createArticlesSQL)
	return
}
