package api

import (
	"net/http"
	"encoding/json"
	util "app/utils"
	"app/models"
	"gopkg.in/go-playground/validator.v9"
)

type LoginInput struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type SignupInput struct {
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type ResendActivationInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ActivateAccountInput struct {
	ActivationCode string `json:"activationCode" validate:"required"`
}

type ForgetPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordInput struct {
	ResetPasswordCode string `json:"resetPasswordCode" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

var Login = func(w http.ResponseWriter, r *http.Request) {
	var errors []string

	input := LoginInput{}
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
	
	// Login in the user
	user := &models.User{}
	resp := user.Login(input.Email, input.Password)
	
	util.Respond(w, resp)
}

var Signup = func(w http.ResponseWriter, r *http.Request) {
	var errors []string

	input := SignupInput{}
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
	user := &models.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Password = input.Password
	
	// Create the account
	resp := user.Create()
	
	util.Respond(w, resp)
}

var ResendActivation = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	
	input := ResendActivationInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		resp := util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors)
		util.Respond(w, resp)
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
	
	user := &models.User{}
	user.Email = input.Email
	resp := user.ResendActivation()
	
	util.Respond(w, resp)
}

var ActivateAccount = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	input := ActivateAccountInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		resp := util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors)
		util.Respond(w, resp)
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
	
	user := &models.User{}
	resp := user.ActivateAccount(input.ActivationCode)
	
	util.Respond(w, resp)
}

var ForgetPassword = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	
	input := ForgetPasswordInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		resp := util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors)
		util.Respond(w, resp)
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
	
	user := &models.User{}
	user.Email = input.Email
	resp := user.ForgetPassword()
	
	util.Respond(w, resp)
}

var ResetPassword = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	input := ResetPasswordInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		resp := util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors)
		util.Respond(w, resp)
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
	
	user := &models.User{}
	resp := user.ResetPassword(input.ResetPasswordCode, input.Password)
	
	util.Respond(w, resp)
}

