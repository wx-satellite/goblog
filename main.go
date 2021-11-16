package main

import (
	_ "github.com/go-sql-driver/mysql"
	"goblog/app/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/logger"
	"net/http"
)

func main() {
	router := bootstrap.SetupRoute()
	bootstrap.SetupDB()
	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.Error(err)
}
