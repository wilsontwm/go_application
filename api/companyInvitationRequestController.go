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
}

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

var JoinCompanyInvitationRequest = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string] interface{}
	userId := r.Context().Value("user") . (uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["id"]) 
	invitationId, _ := uuid.FromString(vars["invitationID"]) 
	
	// Authorization
	if ok := policy.JoinCompanyInvitation(invitationId, userId, companyId); !ok {
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

	invitation := models.GetCompanyInvitationRequest(invitationId)
	resp = invitation.JoinCompanyInvitation(*user)
	
	util.Respond(w, resp)
}