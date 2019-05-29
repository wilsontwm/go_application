package controllers

import (
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
	success, errors := util.GetFlashMessages(w, r)

	data := map[string]interface{}{
		"title": "Signup",
		"appName": appName,
		"errors": errors,
		"success": success,
	}

	err := templates.ExecuteTemplate(w, "signup_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var SignupSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/signup"
	urlStr := restURL.String()

	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace( r.Form.Get("password"))
	retype_password := strings.TrimSpace(r.Form.Get("retype_password"))

	// Check if the retype password matches
	if(password != retype_password) {
		var errors []string
		errors = append(errors, "Retype password does not match.")
		errorJson, _ := json.Marshal(errors)
		errorFlash := []byte(errorJson)
		util.SetFlash(w, "errors", errorFlash)
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

		util.SetErrorSuccessFlash(w, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

var ResendActivationPage = func(w http.ResponseWriter, r *http.Request) {
	success, errors := util.GetFlashMessages(w, r)

	data := map[string]interface{}{
		"title": "Resend Activation",
		"appName": appName,
		"errors": errors,
		"success": success,
	}

	err := templates.ExecuteTemplate(w, "resend_activation_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var ResendActivationSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/resendactivation"
	urlStr := restURL.String()

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

		util.SetErrorSuccessFlash(w, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}


var ForgetPasswordPage = func(w http.ResponseWriter, r *http.Request) {
	success, errors := util.GetFlashMessages(w, r)

	data := map[string]interface{}{
		"title": "Forgotten Password",
		"appName": appName,
		"errors": errors,
		"success": success,
	}

	err := templates.ExecuteTemplate(w, "forget_password_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var ForgetPasswordSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/forgetpassword"
	urlStr := restURL.String()

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

		util.SetErrorSuccessFlash(w, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
