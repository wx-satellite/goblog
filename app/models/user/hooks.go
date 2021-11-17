package user

import (
	"goblog/pkg/password"
	"gorm.io/gorm"
)

// GORM 模型钩子 是在创建、查询、更新、删除等操作之前、之后调用的函数
// 如果返回错误，GORM 将停止后续的操作并回滚事务

// 关于 GORM 的勾子：https://gorm.io/zh_CN/docs/hooks.html

// 创建和更新都会执行的勾子是：save

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.Password = password.Hash(m.Password)
	return
}

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if !password.IsHashed(m.Password) {
		m.Password = password.Hash(m.Password)
	}
	return
}
