package api

import (
	"strings"
	"strconv"
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

type CompanyInvitationInput struct {
	Emails []string `json:"emails"`
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

var InviteToCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 
	company := models.GetCompany(companyId, userId) 

	if user == nil || company == nil  {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	} 

	input := CompanyInvitationInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	emails := util.GetUniqueValues(input.Emails)

	// Create channel to receive the result
	const noOfEmailWorkers int = 10 // Have 10 goroutines to get the emails
	emailJobs := make(chan string, len(emails))
	invitation := make(chan models.CompanyInvitationRequest, len(emails))

	for w := 1; w <= noOfEmailWorkers; w++ {
		go func(id int, emailJobs <-chan string, results chan<- models.CompanyInvitationRequest) {
			for emailInput := range emailJobs {
				result := company.InviteToCompany(emailInput)
				// signal that the routine has completed
				if(result["success"].(bool)) {
					results <- result["data"].(models.CompanyInvitationRequest)
				} else {
					empty := models.CompanyInvitationRequest{}
					results <- empty
				}
			}
		} (w, emailJobs, invitation)
	}

	// Loop through the emails to check if the email can be invited
	for _, email := range emails {
		// Send the email to the email jobs
		emailJobs <- email
	}
	close(emailJobs)

	// Gather the result
	var successfulEmails []interface{}
	var successfulEmailString []string
	for i := 0; i < len(emails) ; i++ {
        successfulEmail := <-invitation
        if successfulEmail.Email != "" {
			successfulEmails = append(successfulEmails, successfulEmail)
			successfulEmailString = append(successfulEmailString, successfulEmail.Email)
		}
	}

	if len(successfulEmails) > 0 {
		emails := strings.Join(successfulEmailString, ", ")
		resp = util.Message(true, http.StatusOK, "You have successfully invited " + emails + " to the company.", errors)
		resp["emails"] = successfulEmails
	} else {
		resp = util.Message(false, http.StatusOK, "No emails have been invited to the company. Please ensure that the emails are not part of the company already or have not been invited before.", errors)
	}

	resp["company"] = company.Name
	
	util.Respond(w, resp)
}

var IndexInviteToCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)

	user := models.GetUser(userId)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 
	company := models.GetCompany(companyId, userId) 

	if user == nil || company == nil  {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
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
	
	resp = company.GetCompanyInvitationList(page)

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