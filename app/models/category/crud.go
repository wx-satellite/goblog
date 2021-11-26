package category

import "goblog/pkg/model"

// GetAll 获取所有分类
func GetAll() (objs []Category, err error) {
	err = model.DB.Find(&objs).Error
	return
}
