package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/user"
)

// ValidateRegistrationForm 验证注册表单
func ValidateRegistrationForm(data user.User) (errs map[string][]string) {
	// 表单规格
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,50", "not_exists:users,name"},
		"email":            []string{"required", "minUTF8:4", "maxUTF8:50", "email", "not_exists:users,email"},
		"password":         []string{"required", "minUTF8:6"},
		"password_confirm": []string{"required"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:格式错误，只允许数字和英文",
			"between:用户名长度需在 3~50 之间",
		},
		"email": []string{
			"required:Email 为必填项",
			"minUTF8:Email 长度需大于 4",
			"maxUTF8:Email 长度需小于 50",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
		"password": []string{
			"required:密码为必填项",
			"minUTF8:长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: ValidatorFlag, // Struct 标签标识符，默认是 json
	}

	// 4. 开始认证
	errs = govalidator.New(opts).ValidateStruct()

	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入密码不匹配！")
	}
	return
}
