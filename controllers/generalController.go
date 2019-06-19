package controllers

import (
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	util "app/utils"
	"github.com/gorilla/mux"
	//"fmt"
)

var HelloPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Homepage",
		"appName": appName,
	}

	err := templates.ExecuteTemplate(w, "welcome_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var LoginPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Login",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "login_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var LoginSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/login"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace( r.Form.Get("password"))

	// Set the input data
	jsonData := map[string]interface{}{
		"email": email,
		"password": password,
	}

	url := r.Header.Get("Referer")
	response, err := util.SendPostRequest(urlStr, jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)

		
		// If login is authenticated
		if(resp["success"].(bool)) {
			userData := resp["data"].(map[string]interface{})
			// Store the user token in the cookie
			SetCookieHandler(w, r, "auth", userData["token"].(string))
			url = "/dashboard/"
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, url, http.StatusFound)
	}
}

var SignupPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Signup",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "signup_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var SignupSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/signup"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace( r.Form.Get("password"))
	retype_password := strings.TrimSpace(r.Form.Get("retype_password"))

	// Check if the retype password matches
	if(password != retype_password) {
		session.AddFlash("Retype password does not match.", "errors")
		session.Save(r, w)
		
		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

		return
	}

	// Set the input data
	jsonData := map[string]interface{}{
		"email": email,
		"password": password,
		"name": name,
	}

	response, err := util.SendPostRequest(urlStr, jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)

		// Send activation email
		if(resp["success"].(bool)) {
			userData := resp["data"].(map[string]interface{})
			activationLink := appURL + "/activate/" + userData["activationCode"].(string)

			subject := appName + " - Activate your account"
			receiver := email
			r := util.NewRequest([]string{receiver}, subject)
			r.Send("views/mail/signup.html", map[string]string{"appName": appName, "username": name, "activationLink": activationLink})
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

var ResendActivationPage = func(w http.ResponseWriter, r *http.Request) {
	
	data := map[string]interface{}{
		"title": "Resend Activation",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "resend_activation_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var ResendActivationSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/resendactivation"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()
	email := strings.TrimSpace(r.Form.Get("email"))

	// Set the input data
	jsonData := map[string]interface{}{
		"email": email,
	}

	response, err := util.SendPostRequest(urlStr, jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)
		
		// Resend activation email
		if(resp["success"].(bool)) {
			userData := resp["data"].(map[string]interface{})
			activationLink := appURL + "/activate/" + userData["activationCode"].(string)

			subject := appName + " - Activate your account"
			receiver := email
			r := util.NewRequest([]string{receiver}, subject)
			r.Send("views/mail/signup.html", map[string]string{"appName": appName, "username": userData["name"].(string), "activationLink": activationLink})
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

var ActivateAccountPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	
	// Set the URL path
	restURL.Path = "/api/activateaccount"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	vars := mux.Vars(r)
	// Set the input data
	jsonData := map[string]interface{}{
		"activationCode": vars["code"],
	}

	response, err := util.SendPostRequest(urlStr, jsonData)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)
		
		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the login page
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

var ForgetPasswordPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Forgotten Password",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "forget_password_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var ForgetPasswordSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/forgetpassword"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()
	email := strings.TrimSpace(r.Form.Get("email"))

	// Set the input data
	jsonData := map[string]interface{}{
		"email": email,
	}

	response, err := util.SendPostRequest(urlStr, jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)
		
		// Resend activation email
		if(resp["success"].(bool)) {
			userData := resp["data"].(map[string]interface{})
			resetLink := appURL + "/resetpassword/" + userData["resetPasswordCode"].(string)

			subject := appName + " - Reset your password"
			receiver := email
			r := util.NewRequest([]string{receiver}, subject)
			r.Send("views/mail/reset_password.html", map[string]string{"appName": appName, "username": userData["name"].(string), "resetLink": resetLink})
		}

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

var ResetPasswordPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Reset Password",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "reset_password_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var ResetPasswordSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/resetpassword"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()	
	vars := mux.Vars(r)
	password := strings.TrimSpace( r.Form.Get("password"))
	retype_password := strings.TrimSpace(r.Form.Get("retype_password"))

	// Check if the retype password matches
	if(password != retype_password) {
		session.AddFlash("Retype password does not match.", "errors")
		session.Save(r, w)
		
		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer") , http.StatusFound)

		return
	}

	// Set the input data
	jsonData := map[string]interface{}{
		"password": password,
		"resetPasswordCode": vars["code"],
	}

	response, err := util.SendPostRequest(urlStr, jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)

		util.SetErrorSuccessFlash(session, w, r, resp)
		// Redirect back to the login page
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

var Custom403Page = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Not authorized",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "custom_403_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var Custom404Page = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Page not found",
		"appName": appName,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "custom_404_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}