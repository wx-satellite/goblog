package middlewares

import "net/http"

// ForceHtml 设置响应头的返回类型为页面形式
func ForceHtml(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
