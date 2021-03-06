package api

import (
	"strings"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	util "app/utils"
	"encoding/json"
	"app/models"
	"app/policy"
	"github.com/satori/go.uuid"
)

type CompanyInvitationInput struct {
	Emails []string `json:"emails"`
	Message string `json:"message"`
}

type CompanyInvitationResponseInput struct {
	IsJoin bool `json:"is_join"`
}

// Send invitation to emails to join company
var InviteToCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.CreateUpdateDeleteCompanyInvitation(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)
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
	message := input.Message

	// Create channel to receive the result
	const noOfEmailWorkers int = 10 // Have 10 goroutines to get the emails
	emailJobs := make(chan string, len(emails))
	invitation := make(chan models.CompanyInvitationRequest, len(emails))

	for w := 1; w <= noOfEmailWorkers; w++ {
		go func(id int, emailJobs <-chan string, results chan<- models.CompanyInvitationRequest) {
			for emailInput := range emailJobs {
				result := company.InviteToCompany(emailInput, message, userId)
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

// Get a list of invited emails to the company
var IndexInviteToCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)

	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.CreateUpdateDeleteCompanyInvitation(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)
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

// Show the company invitation request
var ShowCompanyInvitationRequest = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}

	userId := r.Context().Value("user") . (uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 

	// Authorization
	if ok := policy.ShowCompanyInvitation(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)
	invitationId, _ := uuid.FromString(vars["invitationID"]) 
	company := models.GetCompany(companyId, userId) 

	if user == nil || company == nil  {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	} 

	invitation := &models.CompanyInvitationRequest{}
	resp = invitation.GetInvitation(invitationId, companyId)
	resp["company"] = company

	util.Respond(w, resp)
}

// Delete the company invitation request
var DeleteCompanyInvitationRequest = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 
	
	// Authorization
	if ok := policy.CreateUpdateDeleteCompanyInvitation(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)	
		util.Respond(w, resp)
		return
	}

	user := models.GetUser(userId)
	invitationId, _ := uuid.FromString(vars["invitationID"]) 

	if user == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	} 

	invitation := &models.CompanyInvitationRequest{}
	resp = invitation.GetInvitation(invitationId, companyId)
	if _, ok := resp["data"]; ok {
		data := resp["data"] . (*models.CompanyInvitationRequest)
		resp = data.DeleteInvitation()
	}
	
	util.Respond(w, resp)
}

// User gets all the invitation requests from all the companies
var IndexInvitationFromCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)
	
	user := models.GetUser(userId)

	if user == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)	
		util.Respond(w, resp)
		return
	} 
	
	resp = user.GetCompanyInvitationList()

	util.Respond(w, resp)
}

// User gets the invitation request
var ShowInvitationFromCompany = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)
	
	// Get the ID of the invitation passed in via URL
	vars := mux.Vars(r) 
	invitationId, _ := uuid.FromString(vars["id"]) 
	
	// Authorization
	if ok := policy.ShowInvitationFromCompany(userId, invitationId); !ok {
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

	invitation := models.CompanyInvitationRequest{}
	resp = invitation.GetInvitationFromCompany(invitationId)

	util.Respond(w, resp)
}

// User responds to the company invitation requests, whether to accept or decline invitation request
var RespondCompanyInvitationRequest = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)
	// Get the ID of the invitation passed in via URL
	vars := mux.Vars(r) 
	invitationId, _ := uuid.FromString(vars["id"]) 
	
	// Authorization
	if ok := policy.RespondCompanyInvitation(invitationId, userId); !ok {
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

	input := CompanyInvitationResponseInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	invitation := models.GetCompanyInvitationRequest(invitationId)

	invitationStatus := models.InvitationStatus
	invitationInterface := make([]interface{}, len(invitationStatus))
	for i, v := range invitationStatus {
		invitationInterface[i] = v
	}
	
	invitation.Status = util.IndexOf("Declined", invitationInterface)
	if input.IsJoin == true {
		invitation.Status = util.IndexOf("Joined", invitationInterface)
	}

	resp = invitation.RespondCompanyInvitation(*user)
	
	util.Respond(w, resp)
}