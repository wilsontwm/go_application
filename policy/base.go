package policy

import (
	"app/models"
	"github.com/satori/go.uuid"
)

func IsAdmin(userId, companyId uuid.UUID) bool {
	// Check if user is admin in the company
	user := models.GetUser(userId)
	comp := models.GetCompanyByID(companyId)

	if user == nil || comp == nil {
		return false
	}
	
	return user.IsAdmin(comp)
}