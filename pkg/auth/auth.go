package auth

import (
	"errors"
	"goblog/app/models/user"
	"goblog/pkg/session"
)

// auth 包和业务代码打交道，而不是直接用 session 包，这样子做是为了解耦和隔离
// 如果后续需要更换 session 的驱动（ redis、文件等等，目前是将会话数据存储到cookie中的 ）或者使用其他认证机制
// 就不需要修改业务代码，只需要调整auth包中的相应方法（ auth 包的方法名需要通用，命名上不应该和session包有关系 ）

var UserKey = "uid"

// getUid 获取用户的ID
func getUid() (uid uint64) {
	val := session.Get(UserKey)
	uid, _ = val.(uint64)
	return
}

// User 根据用户ID获取对象
func User() (obj user.User) {
	obj, _ = user.Get(getUid())
	return
}

// Attempt 根据邮箱和密码登陆
func Attempt(email string, password string) (err error) {
	obj, err := user.GetByEmail(email)
	if err != nil {
		return
	}
	if obj.ID <= 0 {
		err = errors.New("账号不存在或者密码错误")
		return
	}

	if !obj.ComparePassword(password) {
		err = errors.New("账号不存在或者密码错误")
		return
	}

	session.Put(UserKey, obj.ID)

	return
}

// Login 登陆指定用户
func Login(obj user.User) {
	session.Put(UserKey, obj.ID)
}

// Logout 登出用户
func Logout() {
	session.Forget(UserKey)
}

// Check 检测是否登陆
func Check() bool {
	return getUid() > 0
}
