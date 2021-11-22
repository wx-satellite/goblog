package user

import (
	"goblog/pkg/model"
	"gorm.io/gorm"
)

// Get 根据用户ID获取对象
func Get(id uint64) (obj User, err error) {
	err = model.DB.Where("id = ?", id).First(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	return
}

// GetByEmail 根据邮箱获取对象
func GetByEmail(email string) (obj User, err error) {
	err = model.DB.Where("email = ?", email).First(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	return
}
