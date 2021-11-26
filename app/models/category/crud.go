package category

import (
	"goblog/pkg/model"
	"gorm.io/gorm"
)

// GetAll 获取所有分类
func GetAll() (objs []Category, err error) {
	err = model.DB.Find(&objs).Error
	return
}

// Find 获取分类实例
func Find(idStr string) (obj Category, err error) {
	err = model.DB.Where("id = ?", idStr).First(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	return
}
