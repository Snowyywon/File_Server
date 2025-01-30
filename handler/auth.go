package handler

import (
	"net/http"
)

// http 拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			// 校验username token是否有效
			if len(username) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w, r)
		})
}
