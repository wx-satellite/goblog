package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/article"
)

// ValidateArticleForm 验证文章表单
func ValidateArticleForm(data article.Article) (errs map[string][]string) {

	// 定义规则
	rules := govalidator.MapData{
		"body":  []string{"required", "minUTF8:10"},
		"title": []string{"required", "minUTF8:3", "maxUTF8:50"},
	}

	// 定义错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"minUTF8:标题长度至少为3",
			"maxUTF8:标题长度不能超过50",
		},
		"body": []string{
			"required:文章内容为必填项",
			"minUTF8:文章内容长度至少为10",
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
