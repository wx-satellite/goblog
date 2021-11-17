package auth

import (
	"errors"
	"goblog/app/models/user"
	"goblog/pkg/session"
)

var UserKey = "uid"

// getUid 获取用户的ID
func getUid() (uid uint64) {
	val := session.Get(UserKey)
	uid, _ = val.(uint64)
	return
}

// User 根据用户ID获取对象
func User() (obj user.User) {
	obj, _ = obj.Get(getUid())
	return
}

// Attempt 根据邮箱和密码登陆
func Attempt(email string, password string) (err error) {
	obj := user.User{}
	obj, err = obj.GetByEmail(email)
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
