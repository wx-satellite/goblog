package article

import (
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/model"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

type Article struct {
	models.BaseModel
	Title  string `gorm:"column:title;type:varchar(50)" valid:"title"`
	Body   string `gorm:"column:body" valid:"body"`
	UserId uint64 `gorm:"not null;index"`
	User   user.User
}

func (article Article) Link() string {
	return route.NameToUrl("articles.show", "id", types.Uint64ToString(article.ID))
}

func (article *Article) Create() (err error) {
	return model.DB.Create(article).Error
}

func (article *Article) Update() (rows int64, err error) {
	// result.RowsAffected 更新的记录数
	// result.Error        更新的错误
	result := model.DB.Save(article) // Save 会保存所有的字段，即使字段是零值，如果没有 primary key 就会新增
	if err = result.Error; err != nil {
		return
	}
	rows = result.RowsAffected
	return
}

func (article *Article) Delete() (rows int64, err error) {
	// 根据主键删除
	// 也可以增加筛选条件： model.DB.Where().Delete()
	result := model.DB.Delete(article)
	if err = result.Error; err != nil {
		return
	}
	rows = result.RowsAffected
	return
}

func (article Article) CreatedAtDate() string {
	return article.CreatedAt.Format("2006-01-02")
}
