package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/view"
	"net/http"
)

// TODO：验证邮箱功能、找回密码功能

type AuthController struct {
}

// Logout 退出登陆
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	flash.Success("您已退出登陆")
	auth.Logout()
	http.Redirect(w, r, "/", http.StatusFound)
}

// Login 登陆页面渲染
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	_ = view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 处理登陆逻辑
func (c *AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	if err := auth.Attempt(email, password); err != nil {
		_ = view.RenderSimple(w, view.D{
			"Email":    email,
			"Password": password,
			"Error":    err,
		}, "auth.login")
		return
	}
	flash.Success("欢迎回来！")
	// 登陆成功
	http.Redirect(w, r, "/", http.StatusFound)
}

// Register 注册页面渲染
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	_ = view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 处理注册逻辑
func (c *AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 在写代码之前，可以先写一些注释( 伪代码 )，帮助整理思路，比如注册的入库逻辑：
	// 	1. 表单验证( 用户提交的数据是不信任的 )
	// 	2. 验证通过，入库
	// 	3. 验证不通过，重新展示创建表单，并提示错误
	userObj := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}
	errs := requests.ValidateRegistrationForm(userObj)
	if len(errs) > 0 {
		_ = view.RenderSimple(w, view.D{
			"User":   userObj,
			"Errors": errs,
		}, "auth.register")
		return
	}

	err := userObj.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "服务器内部错误")
		return
	}

	if userObj.ID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "注册失败，请联系管理员：1453085314@qq.com")
		return
	}

	// 注册成功，自动登陆
	flash.Success("恭喜您注册成功！")
	auth.Login(userObj)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
