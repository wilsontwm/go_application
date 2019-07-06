package api

import (
	"net/http"
	util "app/utils"
	"encoding/json"
	"app/models"
	"gopkg.in/go-playground/validator.v9"
)

type EditProfileInput struct {
	Name string `json:"name" validate:"required"`
	Bio string `json:"bio"`
}

type UploadPictureInput struct {
	ProfilePicture string `json:"profilePicture"`
}

type EditPasswordInput struct {
	Password string `json:"password" validate:"required,min=8,max=16"`
}

// Get the profile information
var GetProfile = func(w http.ResponseWriter, r *http.Request) {
	var errors []string

	userId := r.Context().Value("user") . (uint)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	resp := util.Message(true, http.StatusOK, "Successfully reset the password.", errors)	
	resp["data"] = user
	util.Respond(w, resp)
}

var EditProfile = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uint)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	input := EditProfileInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		util.GetErrorMessages(&errors, err)
		
		resp := util.Message(false, http.StatusUnprocessableEntity, "Validation error", errors)
		util.Respond(w, resp)
		return
	}
	// Save the data into database
	user.Name = input.Name	
	user.Bio = input.Bio
	resp := user.EditProfile()
	
	util.Respond(w, resp)
}

var UploadPicture = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uint)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	input := UploadPictureInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	// Save the data into database
	user.ProfilePicture = input.ProfilePicture	
	resp := user.UploadPicture()
	
	util.Respond(w, resp)
}

var EditPassword = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uint)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	input := EditPasswordInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		util.GetErrorMessages(&errors, err)
		
		resp := util.Message(false, http.StatusUnprocessableEntity, "Validation error", errors)
		util.Respond(w, resp)
		return
	}
	// Save the data into database
	user.Password = input.Password
	resp := user.EditPassword()
	
	util.Respond(w, resp)
}
