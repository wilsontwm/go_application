package api

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	util "app/utils"
	"encoding/json"
	"app/models"
	"app/policy"
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

// Get a list of companies
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

// Create a new company and become admin of the newly created company
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

// Get the detail of the company
var ShowCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)	

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.ShowCompany(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	}

	company := &models.Company{}
	
	resp := company.ShowCompany(companyId, userId)
	
	util.Respond(w, resp)
}

// Update the company
var EditCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.UpdateCompany(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)
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

// Delete the company
var DeleteCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.UpdateCompany(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)
	company := models.GetCompany(companyId, userId) 

	if user == nil || company == nil  {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	} 
	
	resp := company.DeleteCompany()
	
	util.Respond(w, resp)
}

// Get a unique slug for the company name
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

// Get the users in the company
var IndexCompanyUsers = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user") . (uuid.UUID)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.ViewCompanyUsers(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	// Get the page passed in via URL
	pageKeys, ok := r.URL.Query()["page"]
	page := 0 // if page 0, then show all

	if ok && len(pageKeys[0]) >= 1 {
		if _, err := strconv.Atoi(pageKeys[0]); err == nil {
			page, _ = strconv.Atoi(pageKeys[0])
		}
	}	
	
	company := models.GetCompanyByID(companyId)
	resp := company.GetUserList(page)
	
	util.Respond(w, resp)
}