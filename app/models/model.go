package models

import (
	"goblog/pkg/types"
	"time"
)

type BaseModel struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"`
	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

// GetStringID 获取 ID 的字符串形式
func (m BaseModel) GetStringID() string {
	return types.Uint64ToString(m.ID)
}
