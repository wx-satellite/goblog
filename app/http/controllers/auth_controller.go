package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct {
}

// Register 注册页面渲染
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	_ = view.RenderSimple(w, view.D{}, "auth.register")
}

type userForm struct {
	Name            string `valid:"name"`
	Password        string `valid:"password"`
	Email           string `valid:"email"`
	PasswordConfirm string `valid:"password_confirm"` // 确认密码
}

// DoRegister 处理注册逻辑
func (c *AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 在写代码之前，可以先写一些注释( 伪代码 )，帮助整理思路，比如注册的入库逻辑：
	// 	1. 表单验证( 用户提交的数据是不信任的 )
	// 	2. 验证通过，入库
	// 	3. 验证不通过，重新展示创建表单，并提示错误
	userObj := user.User{
		Name:            r.PostFormValue("name"),
		Password:        r.PostFormValue("email"),
		Email:           r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}
	errs := requests.ValidateRegistrationForm(userObj)
	if len(errs) > 0 {
		// 4.1 有错误发生，打印数据
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

	http.Redirect(w, r, "/", http.StatusFound)
	return
}
