package policy

import (
	"github.com/satori/go.uuid"
	"app/models"
)

// Check if the user can create/edit/delete the company invitation request
func CreateUpdateDeleteCompanyInvitation(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}

// Check if the user can see the list of company invitation requests
func ShowCompanyInvitation(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}

// Check if the user can view the invitation from company
func ShowInvitationFromCompany(userId, invitationId uuid.UUID) bool {
	db := models.GetDB()
	defer db.Close()

	// Check if the invitation email is matching
	invitation := models.CompanyInvitationRequest{}
	db.Table("company_invitation_requests").
	Select("company_invitation_requests.*").
	Joins("left join users on users.email = company_invitation_requests.email").
	Where("company_invitation_requests.id = ? AND users.id = ?", invitationId, userId).
	Scan(&invitation)

	return invitation.ID != uuid.Nil
}

// Check if the user can respond to the company invitation request
func RespondCompanyInvitation(invitationId, userId uuid.UUID) bool {
	db := models.GetDB()
	defer db.Close()

	// Check if the invitation email is matching
	invitation := models.CompanyInvitationRequest{}
	db.Table("company_invitation_requests").
	Select("company_invitation_requests.*").
	Joins("left join users on users.email = company_invitation_requests.email").
	Where("company_invitation_requests.id = ? AND company_invitation_requests.status = 0 AND users.id = ?", invitationId, userId).
	Scan(&invitation)

	return invitation.ID != uuid.Nil
}