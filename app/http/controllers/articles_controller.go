package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
	"unicode/utf8"
)

type ArticlesController struct {
}

func (c *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	obj, err := article.Find(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	if obj.ID <= 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "文章不存在")
		return
	}

	_, err = obj.Delete()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
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
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	if obj.ID <= 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "文章不存在")
		return
	}
	_ = view.Render(w, view.D{
		"Title":   obj.Title,
		"Body":    obj.Body,
		"Article": obj,
	}, "articles.edit", "articles.form")
}

func (c *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	obj, err := article.Find(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	if obj.ID <= 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "文章不存在")
		return
	}

	title, body := r.PostFormValue("title"), r.PostFormValue("body")

	errs := validateArticleFormData(title, body)

	// 存在错误
	if len(errs) > 0 {
		err = view.Render(w, view.D{
			"Title":   obj.Title,
			"Body":    obj.Body,
			"Article": obj,
			"Errors":  errs,
		}, "articles.edit", "articles.form")
		return
	}

	// 不存在错误，则更新
	obj.Title = title
	obj.Body = body
	_, err = obj.Update()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	showUrl := route.NameToUrl("articles.show", "id", types.Uint64ToString(obj.ID))
	http.Redirect(w, r, showUrl, http.StatusFound)
	return
}

// validateArticleFormData 表单验证
func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

// Store 文章创建
func (c *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title, body := r.PostFormValue("title"), r.PostFormValue("body")
	errors := validateArticleFormData(title, body)

	// 存在错误时，重新渲染创建表单，并把错误显示出来
	if len(errors) > 0 {
		_ = view.Render(w, view.D{
			"Errors": errors,
			"Title":  title,
			"Body":   body,
		}, "articles.create", "articles.form")
		return
	}

	obj := article.Article{
		Title: title,
		Body:  body,
	}

	// 新增文章
	err := obj.Create()

	logger.Error(err)

	if obj.ID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
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
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}

	_ = view.Render(w, articles, "articles.index")

}

// Show 文章详情页
func (c *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	obj, err := article.Find(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	if obj.ID <= 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "文章不存在")
		return
	}

	_ = view.Render(w, obj, "articles.show")

}
