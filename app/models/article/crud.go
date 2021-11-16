package article

import (
	"goblog/pkg/model"
	"gorm.io/gorm"
)

// Find 根据 ID 获取文章
func Find(idStr string) (obj Article, err error) {
	err = model.DB.Where("id = ?", idStr).First(&obj).Error
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
func GetAll() (objs []Article, err error) {
	err = model.DB.Find(&objs).Error
	return
}
