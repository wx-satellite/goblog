package controllers

import (
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
)

type UserController struct {
	BaseController
}

func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	obj, err := user.Get(types.StringToUint(id))
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	if obj.ID <= 0 {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusNotFound, Message: "用户不存在"})
		return
	}

	articles, err := article.GetAllByUid(id)
	if err != nil {
		c.ResponseError(w, ErrorMessage{HttpCode: http.StatusInternalServerError})
		return
	}
	_ = view.Render(w, view.D{
		"Articles": articles,
	}, "articles.index", "articles._article_meta")
}
