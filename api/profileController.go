package api

import (
	"net/http"
	util "app/utils"
	"encoding/json"
	"app/models"
	"gopkg.in/go-playground/validator.v9"
	"github.com/satori/go.uuid"
	"time"
)

type EditProfileInput struct {
	Name string `json:"name" validate:"required"`
	Phone string `json:"phone"`
	City string `json:"city"`
	Country int `json:"country"`
	Gender int `json:"gender"`
	Birthday *time.Time `json:"birthday"`
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
	countries := models.GetCountries()
	genders := models.GetGenders()
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	resp := util.Message(true, http.StatusOK, "Successfully retrieved the data.", errors)	
	resp["data"] = user
	resp["countries"] = countries
	resp["genders"] = genders
	util.Respond(w, resp)
}

var EditProfile = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

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
	user.Phone = input.Phone
	user.City = input.City
	user.Country = input.Country
	user.Gender = input.Gender
	user.Birthday = input.Birthday
	user.Bio = input.Bio

	if(input.Birthday.IsZero()) {
		user.Birthday = nil
	}
	
	resp := user.EditProfile()
	
	util.Respond(w, resp)
}

var UploadPicture = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

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

var DeletePicture = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	// Save the data into database
	user.ProfilePicture = ""
	resp := user.DeletePicture()
	
	util.Respond(w, resp)
}

var EditPassword = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

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
