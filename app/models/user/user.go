package user

import (
	"goblog/app/models"
	"goblog/pkg/model"
	"gorm.io/gorm"
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

// Get 根据用户ID获取对象
func (m *User) Get(id uint64) (obj User, err error) {
	err = model.DB.Where("id = ?", id).First(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	return
}

// GetByEmail 根据邮箱获取对象
func (m *User) GetByEmail(email string) (obj User, err error) {
	err = model.DB.Where("email = ?", email).First(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	return
}

// ComparePassword 比较密码，相等返回true
func (m *User) ComparePassword(password string) bool {
	return m.Password == password
}
