package api

import (
	"app/models"
	util "app/utils"
	//encoding/json"
	"github.com/satori/go.uuid"
	"net/http"
)

type Payload struct {
	SenderID   string    `json:"sender"`
	ReceiverID string    `json:"receiver"`
	TopicID    uuid.UUID `json:"topic"`
	TopicName  string    `json:"topicname"`
	Msg        string    `json:"message"`
}

var ChatTest = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	util.Respond(w, util.Message(true, http.StatusOK, "TEST", errors))
}

// Create a new chat room and subscribe to topic
var AddChat = func(w http.ResponseWriter, r *http.Request) {
	var errors []string
	var resp map[string]interface{}

	r.ParseForm()
	input := &Payload{
		SenderID:   r.FormValue("sender"),
		ReceiverID: r.FormValue("receiver"),
		TopicName:  r.FormValue("topicname"),
	}

	// Check whether topic already exist
	var topic = models.GetTopicByID(input.TopicID)

	// Create and save the data into database if does not exist
	if topic == nil {
		topic = &models.Topic{
			Name: input.TopicName,
		}
		resp = topic.Create()
	} else {
		resp = util.Message(true, http.StatusOK, "Topic already exists, no new topic will be created", errors)
	}

	util.Respond(w, resp)
}
