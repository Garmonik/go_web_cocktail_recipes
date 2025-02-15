package middleware_auth

import (
	"net/http"
)

// AuthMiddleware check auth_token in cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			// Если токен отсутствует, перенаправляем на /home/
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
