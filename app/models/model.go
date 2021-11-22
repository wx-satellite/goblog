package models

import (
	"goblog/pkg/types"
	"time"
)

// GORM 倾向于约定，而不是配置。默认情况下，GORM 使用 ID 作为主键
// 使用结构体名(OpenUser)的 蛇形复数(open_users) 作为表名，字段名(AppId)的 蛇形(app_id) 作为列名
// 并使用 CreatedAt、UpdatedAt 字段追踪创建、更新时间
type BaseModel struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"`
	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

// GetStringID 获取 ID 的字符串形式
func (m BaseModel) GetStringID() string {
	return types.Uint64ToString(m.ID)
}
