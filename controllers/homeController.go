package controllers

import (
	"net/http"
	util "app/utils"
	"time"
	"io/ioutil"
	"encoding/json"
	"strings"
	//"fmt"
)

var DashboardPage = func(w http.ResponseWriter, r *http.Request) {
	name := ReadCookieHandler(w, r, "name")
	year := time.Now().Year()
	data := map[string]interface{}{
		"title": "Dashboard",
		"appName": appName,
		"name": name,
		"year": year,
	}

	data, err := util.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "dashboard_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var EditProfilePage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	name := ReadCookieHandler(w, r, "name")
	year := time.Now().Year()

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/get"
	urlStr := restURL.String()

	// Get the info for edit profile
	auth := ReadCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})
	response, err := util.SendAuthenticatedRequest(urlStr, "GET", auth, jsonData)
	
	// Check if response is unauthorized
	if response.StatusCode == http.StatusUnauthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseBody, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(responseBody)), &resp)

		data := map[string]interface{}{
			"title": "Edit Profile",
			"appName": appName,
			"name": name,
			"year": year,
			"user": resp["data"].(map[string]interface{}),
		}

		data, err = util.InitializePage(w, r, store, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	
		err = templates.ExecuteTemplate(w, "edit_profile_html", data)
	
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}	
	}
}

var EditProfileSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/edit"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadCookieHandler(w, r, "auth")
	
	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	bio := strings.TrimSpace(r.Form.Get("bio"))
	
	// Set the input data
	jsonData := map[string]interface{}{
		"name": name,
		"bio": bio,
	}

	response, err := util.SendAuthenticatedRequest(urlStr, "POST", auth, jsonData)
	
	// Check if response is unauthorized
	if response.StatusCode == http.StatusUnauthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)		
	
		// Need to reset the cookie that store name
		userData := resp["data"].(map[string]interface{})
		SetCookieHandler(w, r, "name", userData["name"].(string))

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

var EditPasswordSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/edit/password"
	urlStr := restURL.String()

	session, err := util.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadCookieHandler(w, r, "auth")
	
	// Get the input data from the form
	r.ParseForm()
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
		"password": password,
	}

	response, err := util.SendAuthenticatedRequest(urlStr, "POST", auth, jsonData)
	
	// Check if response is unauthorized
	if response.StatusCode == http.StatusUnauthorized {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		
		// Parse it to json data
		json.Unmarshal([]byte(string(data)), &resp)		

		util.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}