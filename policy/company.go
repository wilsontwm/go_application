package policy

import (
	"app/models"
	"github.com/satori/go.uuid"
)

func ShowCompany(userId, companyId uuid.UUID) bool {
	// Check if the user belongs to the company
	company := models.GetCompany(companyId, userId)

	return company != nil
}

func UpdateCompany(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}