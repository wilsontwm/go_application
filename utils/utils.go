package utils

import (	
	"net/http"
	"bytes"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

// Build json message
func Message(success bool, status int, message string, errors []string) (map[string] interface{}) {
	return map[string]interface{} {"success": success, "status": status, "message": message, "errors": errors}
}

// Return json response
func Respond(w http.ResponseWriter, data map[string] interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data["status"].(int))
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

// Get the success / errors flash message from the cookie
func GetFlashMessages(w http.ResponseWriter, r *http.Request) (success string, errors []string) {
	errorsByte, _ := GetFlash(w, r, "errors")
	successByte, _ := GetFlash(w, r, "success")
	json.Unmarshal([]byte(string(errorsByte)), &errors)
	success = string(successByte)

	return
}

// Build the error message
func GetErrorMessages(errors *[]string, err error) {
	for _, errz := range err.(validator.ValidationErrors) {
		// Build the custom errors here
		switch tag := errz.ActualTag(); tag {
			case "required":
				*errors = append(*errors, errz.StructField() + " is required.")
			case "email":
				*errors = append(*errors, errz.StructField() + " is an invalid email address.")
			default:
		}		
	}

	return
}