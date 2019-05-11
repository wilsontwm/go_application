package controllers

import (
	"net/http"
	"fmt"
	"io/ioutil"
	util "app/utils"
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

var SignupSubmit = func(w http.ResponseWriter, r *http.Request) {
	// Set the URL path
	restURL.Path = "/api/signup"
	urlStr := restURL.String()

	// Get the input data
	jsonData := map[string]interface{}{
		"email": "hello@gmail.com",
		"password": "password",
	}

	response, err := util.SendPostRequest(urlStr, jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
