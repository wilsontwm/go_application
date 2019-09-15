package models

import (
	util "app/utils"
	"github.com/satori/go.uuid"
	"net/http"
)

type Message struct {
	Base
	Sender   User      `gorm:"foreignkey:ID"`
	Receiver User      `gorm:"foreignkey:ID"`
	TopicID  uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	Text     string    `json:"text"`
}

// Create new message
func (msg *Message) Create() map[string]interface{} {
	var errors []string

	// Authentication?

	// Save to db
	db := GetDB()
	db.Create(&msg)

	resp := util.Message(true, http.StatusOK, "Successfully saved message.", errors)
	resp["data"] = msg

	return resp
}

// Delete message
func (msg *Message) Delete() map[string]interface{} {
	var errors []string

	// Authentication?

	// Save to db
	db := GetDB()
	db.Delete(&msg)

	resp := util.Message(true, http.StatusOK, "Successfully deleted message.", errors)

	return resp
}
