package policy

import (
	"github.com/satori/go.uuid"
)

func CreateUpdateDeleteCompanyInvitation(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}

func ShowCompanyInvitation(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}