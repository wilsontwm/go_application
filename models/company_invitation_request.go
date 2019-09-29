package models

import (
	util "app/utils"
	"errors"
	"github.com/satori/go.uuid"
	"net/http"
)

type CompanyInvitationRequest struct {
	Base
	CompanyID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	Email     string    `gorm:"not null;primary_key"`
	Message   string
	SenderID  *uuid.UUID `gorm:"type:uuid"`
	Status    int        `gorm:"default:'0'"`
	UserID    *uuid.UUID `gorm:"type:uuid"`
}

type CompanyInvitationRequestOutput struct {
	CompanyInvitationRequest
	CompanyName string
	SenderName  string
	SenderEmail string
	Timestamp   string
}

var InvitationStatus = []string{
	"Awaiting response",
	"Joined",
	"Declined",
}

// Show the company invitation request
func (invitation *CompanyInvitationRequest) GetInvitation(id, companyId uuid.UUID) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	db := GetDB()
	db.Where("id = ? AND company_id = ?", id, companyId).First(&invitation)
	defer db.Close()

	if invitation.ID == uuid.Nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "No available result.", errors)
		return resp
	}

	resp = util.Message(true, http.StatusOK, "The invitation is retrieved.", errors)
	resp["data"] = invitation

	return resp
}

// Delete the company invitation request
func (invitation *CompanyInvitationRequest) DeleteInvitation() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	db := GetDB()
	db.Delete(&invitation)
	defer db.Close()

	resp = util.Message(true, http.StatusOK, "You have successfully deleted the invitation request.", errors)

	return resp
}

// Show the invitation from company
func (invitation *CompanyInvitationRequest) GetInvitationFromCompany(id uuid.UUID) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	db := GetDB()
	db.Where("id = ?", id).First(&invitation)
	defer db.Close()

	if invitation.ID == uuid.Nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "No available result.", errors)
		return resp
	}

	resp = util.Message(true, http.StatusOK, "The invitation is retrieved.", errors)
	resp["data"] = invitation

	return resp
}

// User responds to the invitation request from company
func (invitation *CompanyInvitationRequest) RespondCompanyInvitation(user User) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	if err := invitation.RespondCompanyTransaction(user); err != nil {
		resp = util.Message(false, http.StatusInternalServerError, err.Error(), errors)
		return resp
	}

	resp = util.Message(true, http.StatusOK, "You have successfully responded to the company invitation.", errors)
	resp["data"] = invitation
	resp["company"] = GetCompanyByID(invitation.CompanyID)

	return resp
}

// A transaction of responding to the company invitation request
func (invitation *CompanyInvitationRequest) RespondCompanyTransaction(user User) error {
	db := GetDB()

	defer db.Close()
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Set the user ID
	if invitation.Status == 1 {
		invitation.UserID = &user.ID
	}

	if err := tx.Save(invitation).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Only create the company user if it's a join response
	if invitation.Status == 1 {
		// Get the user role ID in the company
		userRole := Role{}
		db.Where("company_id = ? AND is_admin = ?", invitation.CompanyID, false).First(&userRole)

		if userRole.ID == uuid.Nil {
			tx.Rollback()
			err := errors.New("The user role is not created in the company.")
			return err
		}

		// Associate the user to the company
		companyUser := CompanyUser{
			UserID:    user.ID,
			CompanyID: invitation.CompanyID,
			RoleID:    userRole.ID,
		}

		if err := tx.Where(companyUser).FirstOrCreate(&companyUser).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func GetCompanyInvitationRequest(invitationID uuid.UUID) *CompanyInvitationRequest {
	// Get the invitation by ID
	invitation := &CompanyInvitationRequest{}
	db := GetDB()
	db.Where("id = ?", invitationID).First(invitation)
	defer db.Close()

	if invitation.ID == uuid.Nil {
		return nil
	}

	return invitation
}
