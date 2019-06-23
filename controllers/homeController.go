package controllers

import (
	"net/http"
	util "app/utils"
	"time"
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
	name := ReadCookieHandler(w, r, "name")
	year := time.Now().Year()

	// Set the URL path
	restURL.Path = "/api/dashboard/profile/get"
	urlStr := restURL.String()

	// Get the info for edit profile
	auth := ReadCookieHandler(w, r, "auth")
	jsonData := make(map[string]interface{})
	response, err := util.SendAuthenticatedRequest(urlStr, "GET", auth, jsonData)
	
	// Check if response is forbidden
	if response.StatusCode == http.StatusForbidden {
		http.Redirect(w, r, "/noaccess", http.StatusFound)
	}
	
	data := map[string]interface{}{
		"title": "Dashboard",
		"appName": appName,
		"name": name,
		"year": year,
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