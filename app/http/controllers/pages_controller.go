package controllers

import "net/http"

// PagesController 处理静态页面
type PagesController struct {
}

// Home 首页
func (c *PagesController) Home(w http.ResponseWriter, r *http.Request) {

}

// About 关于
func (c *PagesController) About(w http.ResponseWriter, r *http.Request) {

}

// NotFound 404页面
func (c *PagesController) NotFound(w http.ResponseWriter, r *http.Request) {

}
