package policy

import (
	"app/models"
	"github.com/satori/go.uuid"
)

// Check if the user can see the user profile
func ShowUserProfile(userId uuid.UUID) bool {
	// Check if the user is valid
	user := models.GetUser(userId)

	return user != nil
}