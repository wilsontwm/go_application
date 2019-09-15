package models

import (
	util "app/utils"
	"github.com/satori/go.uuid"
	"net/http"
)

// Topic model structure
// Medium / Topic used in Pub/Sub
type Topic struct {
	Base
	Name      string      `gorm:"not null"`
	UserTopic []UserTopic `gorm:"foreignkey:TopicID"`
	Message   []Message   `gorm:"foreignkey:TopicID"`
}

// Add a new Topic
func (c Topic) Create() map[string]interface{} {
	var errors []string

	// Authentication?

	// Save to db
	db := GetDB()
	db.Create(&c)

	resp := util.Message(true, http.StatusOK, "Topic successfully created.", errors)
	resp["data"] = c

	return resp
}

func (c Topic) Update() map[string]interface{} {
	var errors []string

	// Authentication?

	// Save to db
	db := GetDB()
	db.Model(&c).Update(map[string]interface{}{
		"Name": c.Name,
	})

	defer db.Close()

	resp := util.Message(true, http.StatusOK, "Successfully updated Topic name.", errors)
	resp["data"] = c

	return resp
}

// Get the Topic via ID
func GetTopicByID(id uuid.UUID) *Topic {
	chl := &Topic{}
	db := GetDB()
	db.Table("Topics").Where("id = ?", id).First(chl)
	defer db.Close()

	if chl.ID == uuid.Nil {
		return nil
	}

	return chl
}
