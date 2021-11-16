package models

import "goblog/pkg/types"

type BaseModel struct {
	ID uint64
}

// GetStringID 获取 ID 的字符串形式
func (m BaseModel) GetStringID() string {
	return types.Uint64ToString(m.ID)
}
