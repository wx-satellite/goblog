package controllers

import (
	"fmt"
	"net/http"
)

// PagesController 处理静态页面
type PagesController struct {
}

// Home 首页
func (c *PagesController) Home(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "这是首页")
}

// About 关于
func (c *PagesController) About(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "这是关于页")
}

// NotFound 404页面
func (c *PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "这是404页")
}
