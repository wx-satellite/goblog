package controllers

import (
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/requests"
	"goblog/pkg/flash"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

type CategoryController struct {
	BaseController
}

// Create 分类创建表单渲染
func (*CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	_ = view.Render(w, view.D{}, "categories.create")
}

// Store 分类创建
func (c *CategoryController) Store(w http.ResponseWriter, r *http.Request) {
	obj := category.Category{Name: r.PostFormValue("name")}
	errs := requests.ValidateCategoryForm(obj)
	if len(errs) >= 1 {
		_ = view.Render(w, view.D{
			"Category": obj,
			"Errors":   errs,
		}, "categories.create")
		return
	}
	err := obj.Create()
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError, Message: "创建分类失败，请联系管理员"})
		return
	}
	flash.Success("创建分类成功")
	http.Redirect(w, r, route.NameToUrl("home"), http.StatusFound)
	return
}

// Show 根据分类展示文章
func (c *CategoryController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	obj, err := category.Find(id)

	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}

	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusNotFound, Message: "分类不存在"})
		return
	}

	objs, pagination, err := article.GetAllByCategoryId(obj.ID, r, 15)
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError, Message: "获取文章数据失败"})
		return
	}

	_ = view.Render(w, view.D{
		"Articles":  objs,
		"PagerData": pagination,
	}, "articles.index", "articles._article_meta")

}
