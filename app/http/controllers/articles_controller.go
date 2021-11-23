package controllers

import (
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
)

type ArticlesController struct {
	BaseController
}

func (c *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	obj, err := article.Find(id)
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{
			HttpCode: http.StatusNotFound,
			Message:  "文章不存在",
		})
		return
	}
	if !policies.CanUpdateArticle(obj) {
		c.ResponseForUnauthorized(w, r)
		return
	}

	_, err = obj.Delete()

	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}

	// 重定向到首页
	indexUrl := route.NameToUrl("articles.index")
	http.Redirect(w, r, indexUrl, http.StatusFound)
	return
}

func (c *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	obj, err := article.Find(id)
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusNotFound, Message: "文章不存在"})
		return
	}
	if !policies.CanUpdateArticle(obj) {
		c.ResponseForUnauthorized(w, r)
		return
	}
	_ = view.Render(w, view.D{
		"Article": obj,
	}, "articles.edit", "articles.form")
}

func (c *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	obj, err := article.Find(id)
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusNotFound, Message: "文章不存在"})
		return
	}
	if !policies.CanUpdateArticle(obj) {
		c.ResponseForUnauthorized(w, r)
		return
	}

	obj.Title = r.PostFormValue("title")
	obj.Body = r.PostFormValue("body")

	// 验证表单
	errs := requests.ValidateArticleForm(obj)

	// 存在错误
	if len(errs) > 0 {
		err = view.Render(w, view.D{
			"Article": obj,
			"Errors":  errs,
		}, "articles.edit", "articles.form")
		return
	}

	_, err = obj.Update()
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	showUrl := route.NameToUrl("articles.show", "id", types.Uint64ToString(obj.ID))
	http.Redirect(w, r, showUrl, http.StatusFound)
	return
}

// Store 文章创建
func (c *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	obj := article.Article{
		Title: r.PostFormValue("title"),
		Body:  r.PostFormValue("body"),
	}
	errs := requests.ValidateArticleForm(obj)
	// 存在错误时，重新渲染创建表单，并把错误显示出来
	if len(errs) > 0 {
		_ = view.Render(w, view.D{
			"Errors":  errs,
			"Article": obj,
		}, "articles.create", "articles.form")
		return
	}

	// 新增文章
	err := obj.Create()

	logger.Error(err)

	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}

	// 跳转到文章详情页
	showUrl := route.NameToUrl("articles.show", "id", types.Uint64ToString(obj.ID))
	http.Redirect(w, r, showUrl, http.StatusFound)
	return
}

// Create 文章创建页面
func (c *ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	// 未初始化的map，写是会报错的，读不会报错
	_ = view.Render(w, view.D{}, "articles.create", "articles.form")
}

// Index 文章列表
func (c *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}

	_ = view.Render(w, view.D{
		"Articles": articles,
	}, "articles.index", "articles._article_meta")

}

// Show 文章详情页
func (c *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	obj, err := article.Find(id)
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusNotFound, Message: "文章不存在"})
		return
	}

	// 这里的 obj 是值类型， 不是指针类型，因此它的Link方法不能写成指针接收者
	// 值类型实际是不能调用指针接收者方法的，为什么能调是因为go的语法糖
	// 但是在 go 的 template 中，则没有这种语法糖，因此 Link 需要是值接收者（ 对应 obj 是值类型 ）否则模版会渲染错误
	logger.Error(view.Render(w, view.D{
		"Article":          obj,
		"CanModifyArticle": policies.CanUpdateArticle(obj),
	}, "articles.show", "articles._article_meta"))

}
