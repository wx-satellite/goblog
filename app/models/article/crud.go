package article

import (
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"gorm.io/gorm"
	"net/http"
)

// Find 根据 ID 获取文章
func Find(idStr string) (obj Article, err error) {
	err = model.DB.Preload("User").Where("id = ?", idStr).First(&obj).Error
	// 查询不到的时候会返回 gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	return
}

// GetAll 获取文章列表
// gorm 文档上说明：当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
// 所以 find 找不到记录应该不会返回 ErrRecordNotFound 错误
func GetAll(r *http.Request, perPage int) (objs []Article, viewData pagination.ViewData, err error) {
	// preload 加载关联数据User
	query := model.DB.Model(Article{}).Preload("User").Order("created_at desc")
	_page := pagination.New(r, query, route.NameToUrl("articles.index"), perPage)

	// 获取视图数据
	viewData = _page.Paging()

	// 获取数据
	// 因为 results 的参数是 interface 而不是 []Article ，所以需要设置 Model 为 Article{} 这样子才能找到数据表
	err = _page.Results(&objs)
	return
}

// GetAllByUid 根据用户id获取文章列表
func GetAllByUid(uid string) (objs []Article, err error) {
	err = model.DB.Preload("User").Where("user_id = ?", uid).Find(&objs).Error
	return
}
