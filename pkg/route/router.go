package route

import (
	"github.com/gorilla/mux"
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
	return urlRes.String()
}

// GetRouteVariable 获取url中的路径参数
func GetRouteVariable(param string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[param]
}
