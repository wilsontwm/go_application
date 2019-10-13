package models

import (
	util "app/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
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

type PostOutput struct {
	Post
	Author            User `gorm:"foreignkey:AuthorID"`
	StatusString      string
	UpdatedAtString   string
	ScheduledAtString string
	PublishedAtString string
}

func IndexPost(companyID uuid.UUID, lastID uuid.UUID, lastPublished time.Time, limit int) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}
	fmt.Println(companyID, lastID, lastPublished, limit)
	posts := []PostOutput{}
	db := GetDB()
	if lastID == uuid.Nil || lastPublished.IsZero() {
		db.Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).
			Table("posts").
			Select("posts.*, TO_CHAR(posts.updated_at, '"+util.DateTimeSQLFormat+"') as updated_at_string, TO_CHAR(posts.scheduled_at, '"+util.DateTimeSQLFormat+"') as scheduled_at_string, TO_CHAR(posts.published_at, '"+util.DateTimeSQLFormat+"') as published_at_string").
			Where("company_id = ? AND published_at IS NOT NULL", companyID).
			Order("published_at DESC, ID DESC").
			Limit(limit).
			Find(&posts)
	} else {
		db.Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).
			Table("posts").
			Select("posts.*, TO_CHAR(posts.updated_at, '"+util.DateTimeSQLFormat+"') as updated_at_string, TO_CHAR(posts.scheduled_at, '"+util.DateTimeSQLFormat+"') as scheduled_at_string, TO_CHAR(posts.published_at, '"+util.DateTimeSQLFormat+"') as published_at_string").
			Where("company_id = ? AND ( published_at < ? OR ( published_at = ? AND id < ? ) ) AND published_at IS NOT NULL", companyID, lastPublished, lastPublished, lastID).
			Order("published_at DESC, ID DESC").
			Limit(limit).
			Find(&posts)
	}

	defer db.Close()

	// Post processing of the posts
	var result []PostOutput
	for _, post := range posts {
		post.StatusString = PostStatusArray[post.Status]
		result = append(result, post)
	}

	resp = util.Message(true, http.StatusOK, "You have successfully retrieved "+strconv.Itoa(len(result))+" posts.", errors)
	resp["data"] = result

	return resp
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

// Delete the post
func (post *Post) DeletePost() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	db := GetDB()
	db.Delete(&post)
	defer db.Close()

	resp = util.Message(true, http.StatusOK, "You have successfully deleted the post.", errors)

	return resp
}

// Get the post based on ID
func GetPostByID(id uuid.UUID) *PostOutput {
	post := PostOutput{}
	db := GetDB()
	db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, profile_picture")
	}).
		Table("posts").
		Select("posts.*, TO_CHAR(posts.updated_at, '"+util.DateTimeSQLFormat+"') as updated_at_string, TO_CHAR(posts.scheduled_at, '"+util.DateTimeSQLFormat+"') as scheduled_at_string, TO_CHAR(posts.published_at, '"+util.DateTimeSQLFormat+"') as published_at_string").
		Where("id = ?", id).
		First(&post)
	defer db.Close()

	if post.ID == uuid.Nil {
		return nil
	}

	// Post processing of the posts
	post.StatusString = PostStatusArray[post.Status]

	return &post
}
