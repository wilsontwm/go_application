package policy

import (
	"app/models"
	"github.com/satori/go.uuid"
)

// Check if the user can see the company
func ShowCompany(userId, companyId uuid.UUID) bool {
	// Check if the user belongs to the company
	company := models.GetCompany(companyId, userId)

	return company != nil
}

// Check if the user can update the company
func UpdateCompany(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	return IsAdmin(userId, companyId)
}

// Check if the user can view all the users in the company
func ViewCompanyUsers(userId, companyId uuid.UUID) bool {
	// Check if the user belongs to the company
	company := models.GetCompany(companyId, userId)

	return company != nil
}

// Check if the user can visit the company
func VisitCompany(userId, companyId uuid.UUID) bool {
	// Check if the user belongs to the company
	company := models.GetCompany(companyId, userId)

	return company != nil
}
