package middleware

import (
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"context"
	"app/controllers"
	util "app/utils"
	"app/models"
	"fmt"
)

var CheckAuth = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCookie := controllers.ReadCookieHandler(w, r, "auth")
			
			if authCookie == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}

var IsLoggedIn = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// list of the urls that do not require checking if it's logged in
			toSkip := []string{"/", "/noaccess"}
			requestPath := r.URL.Path // current request path

			// Check if the request need authentication, 
			// If not, then serve the request
			for _, value := range toSkip {				
				if value == requestPath {
					handler.ServeHTTP(w, r)
					return
				}
			}

			authCookie := controllers.ReadCookieHandler(w, r, "auth")
			
			// If cookie has been set, then redirect to dashboard
			if authCookie != "" {
				http.Redirect(w, r, "/dashboard", http.StatusFound)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
var JwtAuthentication = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var errors []string
			// Check for authentication
			response := make(map[string] interface{})
			tokenHeader := r.Header.Get("Authorization")
	
			// If token is missing, then return error code 403 Unauthorized
			if tokenHeader == "" {
				
				response = util.Message(false, http.StatusUnauthorized, "Missing auth token", errors)
				util.Respond(w, response)
				return
			}
	
			// Check if the token format is correct, ie. Bearer {token}
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				response = util.Message(false, http.StatusUnauthorized, "Invalid auth token format.", errors)
				util.Respond(w, response)
				return
			}
	
			tokenPart := splitted[1] // Grab the second part
			tk := &models.Token{}
	
			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})
	
			if err != nil {
				response = util.Message(false, http.StatusUnauthorized, "Invalid auth token format.", errors)
				util.Respond(w, response)
				return
			}
	
			if !token.Valid {
				response = util.Message(false, http.StatusUnauthorized, "Token is not valid.", errors)
				util.Respond(w, response)
				return
			}
	
			// Everything is authenticated
			fmt.Sprintf("Login User: %s", tk.UserId)

			// Set the user ID in the context
			ctx := context.WithValue(r.Context(), "user", tk.UserId)
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		})
	}
}
