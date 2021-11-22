package policies

import (
	"goblog/app/models/article"
	"goblog/pkg/auth"
)

// CanUpdateArticle 是否可以修改文章
func CanUpdateArticle(obj article.Article) (can bool) {
	return auth.User().ID == obj.UserId
}
