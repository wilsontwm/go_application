package middleware

import (
	"fmt"
	"time"
	"net/http"
	"github.com/gorilla/mux"
)

func LogTime() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("This request was sent at: ", time.Now())
			handler.ServeHTTP(w, r)
		})
	}
}

func Second() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("This second request was sent at: ", time.Now())
			handler.ServeHTTP(w, r)
		})
	}
}

