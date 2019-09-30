package policy

import (
	"app/models"
	"github.com/satori/go.uuid"
)

// Check if the user can create post
func CreatePost(userId, companyId uuid.UUID) bool {
	// Check if the user belongs to the company
	company := models.GetCompany(companyId, userId)

	return company != nil
}
