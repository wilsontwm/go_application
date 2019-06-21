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