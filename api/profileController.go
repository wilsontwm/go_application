package api

import (
	"net/http"
	util "app/utils"
)

// Get the profile information
var GetProfile = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	resp := util.Message(true, http.StatusOK, "Successfully reset the password.", errors)	
	util.Respond(w, resp)
}
