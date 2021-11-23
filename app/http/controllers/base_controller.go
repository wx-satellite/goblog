package controllers

import (
	"fmt"
	"goblog/pkg/flash"
	"net/http"
)

// BaseController 基类控制器
type BaseController struct {
}

type ErrorMessage struct {
	HttpCode int
	Message  string
}

// ResponseError 返回错误
func (c *BaseController) ResponseError(w http.ResponseWriter, message ErrorMessage) {
	w.WriteHeader(message.HttpCode)
	if message.Message == "" {
		message.Message = "服务器内部错误"
	}
	_, _ = fmt.Fprint(w, message.Message)
	return
}

// ResponseForUnauthorized 处理未授权访问
func (c *BaseController) ResponseForUnauthorized(w http.ResponseWriter, r *http.Request) {
	flash.Warning("未授权访问")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
