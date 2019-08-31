package policy

import (
	"github.com/satori/go.uuid"
	"app/models"
)

func CreateUpdateDeleteCompanyInvitation(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}

func ShowCompanyInvitation(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}

func JoinCompanyInvitation(invitationId, userId, companyId uuid.UUID) bool {
	db := models.GetDB()
	defer db.Close()

	// Check if the invitation email is matching
	invitation := models.CompanyInvitationRequest{}
	db.Raw("SELECT INV.* FROM company_invitation_requests INV JOIN users USERS ON USERS.email = INV.email WHERE INV.id = ? AND INV.company_id = ? AND USERS.id = ? AND INV.deleted_at is NULL", invitationId, companyId, userId).Scan(&invitation)

	return invitation.ID != uuid.Nil
}