package category

import (
	"goblog/app/models"
	"goblog/pkg/model"
)

type Category struct {
	models.BaseModel
	Name string `gorm:"column:name;type:varchar(50);not null;default:''" valid:"name"`
}

// Create 创建模型
func (m *Category) Create() (err error) {
	err = model.DB.Create(m).Error
	return
}
