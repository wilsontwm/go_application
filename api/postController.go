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
	"time"
)

type PostInput struct {
	Title       string     `json:"title" validate:"required"`
	Content     string     `json:"content" validate:"required"`
	Status      int        `json:"status" validate:"min=0,max=2"`
	ScheduledAt *time.Time `json:"scheduled_at"`
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
