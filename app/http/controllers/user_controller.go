package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
)

type UserController struct {
}

func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	obj, err := user.Get(types.StringToUint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	if obj.ID <= 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "用户不存在")
	}

	articles, err := article.GetAllByUid(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}
	_ = view.Render(w, view.D{
		"Articles": articles,
	}, "articles.index", "articles._article_meta")
}
