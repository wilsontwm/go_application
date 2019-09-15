package models

import (
	util "app/utils"
	"github.com/satori/go.uuid"
	"net/http"
)

// User connected Topics
type UserTopic struct {
	TopicID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;primary_key"`
}

// Add new User Topic
func (utp UserTopic) Create() map[string]interface{} {
	var errors []string

	// Authentication?

	// Save to db
	db := GetDB()
	db.Create(&utp)

	resp := util.Message(true, http.StatusOK, "Chat room successfully created.", errors)
	resp["data"] = utp

	return resp
}

// Remove user topic
func (utp UserTopic) Delete() map[string]interface{} {
	var errors []string

	// Authentication?

	// Delete entry in db
	db := GetDB()
	usrID := utp.UserID
	usr := GetUser(usrID)

	db.Delete(&utp)

	resp := util.Message(true, http.StatusOK, "Successfully removed user "+usr.Name+".", errors)
	resp["data"] = utp

	return resp
}

// Get user topic
func GetUserTopic(topicID string, userID uuid.UUID) *UserTopic {
	// Authentication?

	// Get entry
	usrTopic := &UserTopic{}
	db := GetDB()
	db.Table("user_topics").Where("topic_id = ? and user_id = ?", topicID, userID).First(&usrTopic)

	return usrTopic
}
