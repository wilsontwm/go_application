package api

import (
	"net/http"
	util "app/utils"
	"encoding/json"
	"app/models"
	"gopkg.in/go-playground/validator.v9"
	"github.com/satori/go.uuid"
)

type CompanyInput struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
	Description string `json:"description"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Fax string `json:"fax"`
	Address string `json:"address"`
}

var IndexCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}
	
	resp := user.IndexCompany()
	
	util.Respond(w, resp)
}

var CreateCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(true, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	input := CompanyInput{}
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

	company := models.Company{
		Name: input.Name,
		Slug: input.Slug,
		Description: input.Description,
		Email: input.Email,
		Phone: input.Phone,
		Fax: input.Fax,
		Address: input.Address,
	}
	
	resp := user.CreateCompany(&company)
	
	util.Respond(w, resp)
}
