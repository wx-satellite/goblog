package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "请求路径："+r.URL.Path)
	})

	http.ListenAndServe(":3000", nil)
}
