package api

import (
	"app/models"
	"app/policy"
	util "app/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostInput struct {
	Title       string     `json:"title" validate:"required"`
	Content     string     `json:"content" validate:"required"`
	Status      int        `json:"status" validate:"min=0,max=2"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

// Get all post
var IndexPost = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user").(uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["companyId"])

	// Authorization
	if ok := policy.IndexPost(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)
		util.Respond(w, resp)
		return
	}

	// Get the published datetime of the last post viewed
	lastPublishedQuery, ok := r.URL.Query()["lastPublished"]
	lastPublished := time.Time{}

	if ok && len(lastPublishedQuery[0]) >= 1 {
		r := strings.NewReplacer(" ", "+")
		lastPublishedString := r.Replace(lastPublishedQuery[0])
		lastPublished, _ = time.Parse(time.RFC3339, lastPublishedString)
	}

	// Get the updated datetime of the last post viewed
	lastUpdatedQuery, ok := r.URL.Query()["lastUpdated"]
	lastUpdated := time.Time{}

	if ok && len(lastUpdatedQuery[0]) >= 1 {
		r := strings.NewReplacer(" ", "+")
		lastUpdatedString := r.Replace(lastUpdatedQuery[0])
		lastUpdated, _ = time.Parse(time.RFC3339, lastUpdatedString)
	}

	// Get the post ID of the last post viewed
	lastIDQuery, ok := r.URL.Query()["lastID"]
	lastID := uuid.Nil
	if ok && len(lastIDQuery[0]) >= 1 {
		lastID, _ = uuid.FromString(lastIDQuery[0])
	}

	// Get the limit of the posts in a batch
	limitQuery, ok := r.URL.Query()["limit"]
	limit := 10 // Default to be 10
	if ok && len(limitQuery[0]) >= 1 {
		if i, err := strconv.ParseInt(limitQuery[0], 10, 32); err == nil {
			limit = int(i)
		}
	}

	// Get the author ID
	authorIDQuery, ok := r.URL.Query()["author"]
	authorID := uuid.Nil
	if ok && len(authorIDQuery[0]) >= 1 {
		authorID, _ = uuid.FromString(authorIDQuery[0])
	}

	// Get the status of the posts
	postStatusQuery, ok := r.URL.Query()["status"]
	postStatus := 2
	if ok && len(postStatusQuery[0]) >= 1 {
		if i, err := strconv.ParseInt(postStatusQuery[0], 10, 32); err == nil {
			pt := int(i)
			if pt >= 0 && pt < len(models.PostStatusArray) {
				postStatus = pt
			}
		}
	}

	resp := models.IndexPost(companyId, lastID, lastPublished, lastUpdated, authorID, postStatus, limit)
	util.Respond(w, resp)
}

// Create a new post
var CreatePost = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user").(uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["companyId"])

	// Authorization
	if ok := policy.CreatePost(userId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)
		util.Respond(w, resp)
		return
	}

	input := PostInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		util.GetErrorMessages(&errors, err)

		resp := util.Message(false, http.StatusUnprocessableEntity, "Validation error", errors)
		util.Respond(w, resp)
		return
	}

	var scheduledAt *time.Time
	var publishedAt *time.Time
	var now = time.Now()
	if models.PostStatusArray[input.Status] == models.PostStatusScheduled {
		scheduledAt = input.ScheduledAt
	} else if models.PostStatusArray[input.Status] == models.PostStatusPublished {
		publishedAt = &now
	}

	post := &models.Post{
		Title:       input.Title,
		Content:     input.Content,
		AuthorID:    userId,
		CompanyID:   companyId,
		Status:      input.Status,
		ScheduledAt: scheduledAt,
		PublishedAt: publishedAt,
	}

	resp := post.CreatePost()

	util.Respond(w, resp)
}

// Get existing post
var ShowPost = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user").(uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["companyId"])
	postId, _ := uuid.FromString(vars["id"])

	// Authorization
	if ok := policy.ShowPost(userId, postId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)
		util.Respond(w, resp)
		return
	}

	post := models.GetPostByID(postId)
	ps := models.GetPostStatus()

	if post == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)
		util.Respond(w, resp)
		return
	}

	resp := util.Message(true, http.StatusOK, "You have successfully retrieved the post.", errors)
	resp["data"] = post
	resp["postStatus"] = ps

	util.Respond(w, resp)
}

// Edit existing post
var EditPost = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user").(uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["companyId"])
	postId, _ := uuid.FromString(vars["id"])

	// Authorization
	if ok := policy.UpdateDeletePost(userId, postId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)
		util.Respond(w, resp)
		return
	}

	post := models.GetPostByID(postId)

	if post == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)
		util.Respond(w, resp)
		return
	}

	input := PostInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errors = append(errors, err.Error())
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body", errors))
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		util.GetErrorMessages(&errors, err)

		resp := util.Message(false, http.StatusUnprocessableEntity, "Validation error", errors)
		util.Respond(w, resp)
		return
	}

	var scheduledAt *time.Time
	var publishedAt *time.Time
	var now = time.Now()
	if models.PostStatusArray[input.Status] == models.PostStatusScheduled {
		scheduledAt = input.ScheduledAt
	} else if models.PostStatusArray[input.Status] == models.PostStatusPublished {
		publishedAt = &now
	}

	post.Title = input.Title
	post.Content = input.Content
	post.Status = input.Status
	post.ScheduledAt = scheduledAt
	post.PublishedAt = publishedAt

	resp := post.EditPost()

	util.Respond(w, resp)
}

// Delete existing post
var DeletePost = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	userId := r.Context().Value("user").(uuid.UUID)
	// Get the ID of the company passed in via URL
	vars := mux.Vars(r)
	companyId, _ := uuid.FromString(vars["companyId"])
	postId, _ := uuid.FromString(vars["id"])

	// Authorization
	if ok := policy.UpdateDeletePost(userId, postId, companyId); !ok {
		resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)
		util.Respond(w, resp)
		return
	}

	post := models.GetPostByID(postId)

	if post == nil {
		resp := util.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.", errors)
		util.Respond(w, resp)
		return
	}

	resp := post.DeletePost()

	util.Respond(w, resp)
}
