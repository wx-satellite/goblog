package main

import (
	_ "github.com/go-sql-driver/mysql"
	"goblog/app/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	"goblog/pkg/logger"
	"net/http"
)

func init() {
	// 初始化配置信息
	// 方法内容为空，只是为了触发 goblog/config 包的 init 方法，其实也可以使用该方法：_ "goblog/config"
	config.Initialize()
}
func main() {
	router := bootstrap.SetupRoute()
	bootstrap.SetupDB()
	bootstrap.SetUpStore()
	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.Error(err)
}
