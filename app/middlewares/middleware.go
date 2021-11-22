package middlewares

import "net/http"

type HttpHandleFunc func(w http.ResponseWriter, r *http.Request)
