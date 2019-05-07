package controllers

import (
	"net/http"
)

var HelloPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Homepage",
		"appName": appName,
	}

	err := templates.ExecuteTemplate(w, "home_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var LoginPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Login",
		"appName": appName,
	}

	err := templates.ExecuteTemplate(w, "login_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var SignupPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Signup",
		"appName": appName,
	}

	err := templates.ExecuteTemplate(w, "signup_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
