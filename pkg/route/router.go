package route

import (
	"github.com/gorilla/mux"
	"goblog/pkg/config"
	"net/http"
)

var Router *mux.Router

// Initialize 初始化路由器
func Initialize() {
	Router = mux.NewRouter()
}

// NameToUrl 根据路由的名称获取url
func NameToUrl(name string, params ...string) string {
	urlRes, _ := Router.Get(name).URL(params...)
	return config.GetString("app.url") + urlRes.String()
}

// GetRouteVariable 获取url中的路径参数
func GetRouteVariable(param string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[param]
}
