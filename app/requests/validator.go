package requests

import (
	"errors"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"goblog/pkg/model"
	"strconv"
	"strings"
	"unicode/utf8"
)

// 表单验证：
//	比较知名的有 asaskevich/govalidator 和 thedevsaddam/govalidator ，两个都值得使用
//	后者借鉴了 Laravel，比较简单易用，本项目将采用此包。

// 验证的标签名称
var ValidatorFlag = "valid"

func init() {

	// 自定义规则，在使用之前需要先注册
	// not_exists:users,email
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) (err error) {

		// 解析规格，取出表名和字段名
		splits := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")
		dbName := splits[0]
		dbField := splits[1]

		// 断言，获取用户提交的值
		val, _ := value.(string)

		var count int64
		model.DB.Table(dbName).Where(fmt.Sprintf("%s = ?", dbField), val).Count(&count)
		// 存在
		if count >= 1 {
			err = fmt.Errorf("%v 已被占用", value)
			if message != "" {
				err = errors.New(message)
			}
		}
		return
	})

	// 自带的 min 和 max 对于中文字符是按照 3 个字节计算的
	govalidator.AddCustomRule("minUTF8", func(field string, rule string, message string, value interface{}) (err error) {
		splits := strings.Split(strings.TrimPrefix(rule, "min:"), ",")
		length, _ := strconv.ParseInt(splits[0], 10, 64)
		valStr, _ := value.(string)
		if int64(utf8.RuneCountInString(valStr)) < length {
			err = errors.New(message)
		}
		return
	})

	govalidator.AddCustomRule("maxUTF8", func(field string, rule string, message string, value interface{}) (err error) {
		splits := strings.Split(strings.TrimPrefix(rule, "max:"), ",")
		length, _ := strconv.ParseInt(splits[0], 10, 64)
		valStr, _ := value.(string)
		if int64(utf8.RuneCountInString(valStr)) > length {
			err = errors.New(message)
		}
		return
	})
}
