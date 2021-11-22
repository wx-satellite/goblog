package session

import (
	"github.com/gorilla/sessions"
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"
)

var CookieKey = config.GetString("session.session_name")

// Store 存储器，这里使用全cookie
var Store = sessions.NewCookieStore([]byte(config.GetString("app.key")))

// NewCookieStore 函数的注释表明，认证密钥推荐使用32位或者64位，但是加密密钥必须是 16、24或32字节
// 对应 AES-128、AES-192或AES-256模式。加密密钥可以不用传递。

// Session 当前会话
var Session *sessions.Session

// Request 用以获取会话
var Request *http.Request

// Response 用以写入会话
var Response http.ResponseWriter

// StartSession 初始化会话，在中间件中调用
func StartSession(w http.ResponseWriter, r *http.Request) {
	var err error

	// Store.Get() 的第二个参数是 Cookie 的名称
	// gorilla/sessions 支持多会话，本项目我们只使用单一会话即可
	Session, err = Store.Get(r, CookieKey)
	logger.Error(err)
	Request = r
	Response = w
}

// Put 写入键值对应的会话数据
func Put(key string, value interface{}) {
	Session.Values[key] = value
	Save()
}

// Get 获取会话数据，获取数据时请做类型检测
func Get(key string) interface{} {
	return Session.Values[key]
}

// Forget 删除某个会话项
func Forget(key string) {
	delete(Session.Values, key)
	Save()
}

// Flush 删除当前会话
func Flush() {
	Session.Options.MaxAge = -1
	Save()
}

// Save 保持会话
func Save() {
	// HttpOnly 为 true 则 js 获取不到 cookie，这样能有效的防止XSS攻击，窃取cookie内容
	// Secure 一般在 https 中使用
	//Session.Options.Secure = true
	//Session.Options.HttpOnly = true

	// 有效期默认是30天，修改可以通过如下方式
	//Session.Options.MaxAge = -1

	err := Session.Save(Request, Response)
	logger.Error(err)
}
