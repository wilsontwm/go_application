package api

import (
	"log"
	"net/http"
	"encoding/json"
	util "app/utils"
	"gopkg.in/go-playground/validator.v9"
)

type SignupInput struct {
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

var Signup = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	input := SignupInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Print("Error decoding request body", err)
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		for _, errz := range err.(validator.ValidationErrors) {
			// Build the custom errors here
			errors = append(errors, errz.Field())
		}

		resp := util.Message(false, http.StatusUnprocessableEntity, "Validation error", errors)
		util.Respond(w, resp)
		return
}

	// Remove the password to be outputted
	input.Password = ""

	resp := util.Message(true, http.StatusOK, "You have successfully signed up.", errors)
	resp["data"] = input
	util.Respond(w, resp)
}
