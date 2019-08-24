package api

import (
	"net/http"
	"github.com/gorilla/mux"
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
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
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
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
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

var ShowCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	company := &models.Company{}
	
	resp := company.ShowCompany(companyId, userId)
	
	util.Respond(w, resp)
}

var EditCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 
	company := models.GetCompany(companyId, userId) 

	if user == nil || company == nil  {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
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

	company.Name = input.Name
	company.Slug = input.Slug
	company.Description = input.Description
	company.Email = input.Email
	company.Phone = input.Phone
	company.Fax = input.Fax
	company.Address = input.Address
	
	resp := company.EditCompany()
	
	util.Respond(w, resp)
}

var DeleteCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 
	company := models.GetCompany(companyId, userId) 

	if user == nil || company == nil  {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	} 
	
	resp := company.DeleteCompany()
	
	util.Respond(w, resp)
}

var GetUniqueSlug = func(w http.ResponseWriter, r *http.Request) {
	compQuery, ok := r.URL.Query()["comp"]
	companyId := uuid.Nil
	if ok && len(compQuery[0]) >= 1 {
		companyId, _ = uuid.FromString(compQuery[0])
	}

	slugQuery, ok := r.URL.Query()["slug"]
	slug := ""
	if ok && len(slugQuery[0]) >= 1 {
		slug = slugQuery[0]
	}

	resp := models.GetUniqueSlug(companyId, slug)
	
	util.Respond(w, resp)
}