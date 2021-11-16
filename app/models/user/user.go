package user

import (
	"goblog/app/models"
	"goblog/pkg/model"
)

type User struct {
	models.BaseModel
	// GORM 默认会将键小写化作为字段名称，并且默认是允许 NULL 的
	Name     string `gorm:"type:varchar(50);not null;default:'';unique" valid:"name"`
	Email    string `gorm:"type:varchar(50);default:'';not null;unique;" valid:"email"`
	Password string `gorm:"type:char(32);default:'';not null" valid:"password"`

	// gorm:"-"  设置 GORM 在读写时略过此字段
	PasswordConfirm string ` gorm:"-" valid:"password_confirm"`
}

func (m *User) Create() (err error) {
	return model.DB.Create(m).Error
}
