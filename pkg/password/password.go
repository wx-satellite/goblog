package password

import (
	"goblog/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// Golang 官方包  crypto/bcrypt 是非常棒的密码加密解决方案，同一个字符串，例如说 abc123456
// 每一次加密出来的结果都不一样，也就是不可逆。
// 因此比较明文和密文是否一样，需要使用该包的 CompareHashAndPassword 方法，不能
// 对明文进行 GenerateFromPassword 加密和密文比较，因此每一次加密的结果都是不一样的

// 如何使用 bcrypt 包？
// 首先我们将 bcrypt 封装到新建的 password 包中，并统一在这个包里做好错误处理；
// 利用 GORM 提供的模型钩子，在创建和更新时对密码进行加密；
// 登录时拿用户提交的明文密码与数据库里的加密过的密码进行配对。

// Hash 使用 bcrypt 对密码进行加密
func Hash(password string) string {
	// GenerateFromPassword 的第二个参数是 cost 值。建议大于 12，数值越大耗费时间越长
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.Error(err)

	return string(bytes)
}

// CheckHash 对比明文密码和数据库的哈希值
func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	logger.Error(err)
	return err == nil
}

// IsHashed 判断字符串是否是哈希过的数据
func IsHashed(str string) bool {
	// bcrypt 加密后的长度等于 60
	return len(str) == 60
}
