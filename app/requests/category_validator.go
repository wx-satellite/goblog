package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/category"
)

// ValidateCategoryForm 验证分类创建表单
func ValidateCategoryForm(obj category.Category) (errs map[string][]string) {
	rules := govalidator.MapData{
		"name": []string{"required", "minUTF8:3", "maxUTF8:50", "not_exists:category,name"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:名称必须填写",
			"minUTF8:长度至少为3",
			"maxUTF8:长度不能超过50",
		},
	}

	options := govalidator.Options{
		Data:          &obj,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}
	return govalidator.New(options).ValidateStruct()
}
