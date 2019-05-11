package api

import (
	"log"
	"net/http"
	"encoding/json"
	util "app/utils"
)

type Input struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

var Signup = func(w http.ResponseWriter, r *http.Request) {
	input := Input{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Print("Error decoding request body", err)
		util.Respond(w, util.Message(false, http.StatusInternalServerError, "Error decoding request body"))
		return
	}

	// Remove the password to be outputted
	input.Password = ""

	resp := util.Message(true, http.StatusOK, "Success")
	resp["data"] = input
	util.Respond(w, resp)
}
