package models

import (
	util "app/utils"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

const PostStatusDraft = "Draft"
const PostStatusScheduled = "Scheduled"
const PostStatusPublished = "Published"

var PostStatusArray = [...]string{
	PostStatusDraft,
	PostStatusScheduled,
	PostStatusPublished,
}

type Post struct {
	Base
	Title       string    `gorm:"not null;"`
	Content     string    `gorm:"not null;"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null"`
	Status      int       `gorm:"not null;default:0"`
	ScheduledAt *time.Time
	PublishedAt *time.Time
}

func (post *Post) Validate() (map[string]interface{}, bool) {
	var errors []string
	var resp map[string]interface{}

	// Check if the status is scheduled, then scheduled at must be set
	if PostStatusArray[post.Status] == PostStatusScheduled && post.ScheduledAt == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Schedule datetime must be set.", errors)
		return resp, false
	}

	// Check if the scheduled at is at least 15 mins later
	currentTime := time.Now().Local().Add(time.Minute * time.Duration(15))
	if PostStatusArray[post.Status] == PostStatusScheduled && post.ScheduledAt.Before(currentTime) {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Schedule datetime must be at least 15 minutes later.", errors)
		return resp, false
	}

	resp = util.Message(true, http.StatusOK, "Input has been validated.", errors)
	return resp, true
}

// Create the post
func (post *Post) CreatePost() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	// Validate the input first
	if resp, ok := post.Validate(); !ok {
		return resp
	}

	if err := CreateUpdatePostTransaction(post); err != nil {
		resp = util.Message(false, http.StatusInternalServerError, err.Error(), errors)
		return resp
	}

	resp = util.Message(true, http.StatusOK, "You have successfully created a post.", errors)
	resp["data"] = post

	return resp
}

// Edit the post
func (post *Post) EditPost() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	// Validate the input first
	if resp, ok := post.Validate(); !ok {
		return resp
	}

	if err := CreateUpdatePostTransaction(post); err != nil {
		resp = util.Message(false, http.StatusInternalServerError, err.Error(), errors)
		return resp
	}

	resp = util.Message(true, http.StatusOK, "You have successfully updated the post.", errors)
	resp["data"] = post

	return resp
}

// The database transaction to create post
func CreateUpdatePostTransaction(post *Post) error {
	db := GetDB()

	defer db.Close()
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		return err
	}

	// TODO: Add post picture & tags

	return tx.Commit().Error
}

// Get the post based on ID
func GetPostByID(id uuid.UUID) *Post {
	post := &Post{}
	db := GetDB()
	db.Table("posts").Where("id = ?", id).First(post)
	defer db.Close()

	if post.ID == uuid.Nil {
		return nil
	}

	return post
}
