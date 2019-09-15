package api

import (
	"net/http"
	"github.com/gorilla/mux"
	util "app/utils"
	"app/models"
	"app/policy"
	"github.com/satori/go.uuid"
)

// Get the user profile information
var GetUserProfile = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	countries := models.GetCountries()
	genders := models.GetGenders()
	userId := r.Context().Value("user") . (uuid.UUID)

	// Authorization
	if ok := policy.ShowUserProfile(userId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	targetUserId, _ := uuid.FromString(vars["id"]) 

	user := models.GetUser(targetUserId)

	if user == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "No result available.", errors)	
		util.Respond(w, resp)
		return
	}

	resp := util.Message(true, http.StatusOK, "Successfully retrieved the data.", errors)	
	resp["data"] = user
	resp["countries"] = countries
	resp["genders"] = genders
	util.Respond(w, resp)
}
