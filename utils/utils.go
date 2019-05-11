package utils

import (	
	"net/http"
	"bytes"
	"encoding/json"
)

// Build json message
func Message(status bool, message string) (map[string] interface{}) {
	return map[string]interface{} {"status": status, "message": message}
}

// Return json response
func Respond(w http.ResponseWriter, data map[string] interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Send a post request to the url
func SendPostRequest(url string, data map[string]interface{}) (response *http.Response, err error) {
	requestBody, err := json.Marshal(data)
	
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}

	response, err = client.Do(request)

	return
}