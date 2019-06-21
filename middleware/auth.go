package middleware

import (
	"net/http"
	"app/controllers"
	"github.com/gorilla/mux"
)

func CheckAuth() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCookie := controllers.ReadCookieHandler(w, r, "auth")
			
			if authCookie == "" {
				http.Redirect(w, r, "/noaccess", http.StatusFound)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}