package models

import (
	util "app/utils"
	"github.com/satori/go.uuid"
	"net/http"
)

type CompanyInvitationRequest struct {
	Base
	CompanyID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	Email string `gorm:"not null;primary_key"`
	UserID *uuid.UUID `gorm:"type:uuid"`
	Status int `gorm:"default:'0'"`
}

func (invitation *CompanyInvitationRequest) GetInvitation(id, companyId uuid.UUID) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
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

func (invitation *CompanyInvitationRequest) DeleteInvitation() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
	db := GetDB()
	db.Delete(&invitation)
	defer db.Close()

	resp = util.Message(true, http.StatusOK, "You have successfully deleted the invitation request.", errors)

	return resp
}