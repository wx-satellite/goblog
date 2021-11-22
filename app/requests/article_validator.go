package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/article"
)

// ValidateArticleForm 验证文章表单
func ValidateArticleForm(data article.Article) (errs map[string][]string) {

	// 定义规则
	rules := govalidator.MapData{
		"body":  []string{"required", "min:10"},
		"title": []string{"required", "min:3", "max:50"},
	}

	// 定义错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min:标题长度至少为3",
			"max:标题长度不能超过50",
		},
		"body": []string{
			"required:文章内容为必填项",
			"min:文章内容长度至少为10",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: ValidatorFlag,
	}

	return govalidator.New(opts).ValidateStruct()
}
