package utils

import (	
	"net/http"
	"reflect"
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

// Build the error message
func GetErrorMessages(errors *[]string, err error) {
	for _, errz := range err.(validator.ValidationErrors) {
		// Build the custom errors here
		switch tag := errz.ActualTag(); tag {
			case "required":
				*errors = append(*errors, errz.StructField() + " is required.")
			case "email":
				*errors = append(*errors, errz.StructField() + " is an invalid email address.")
			case "min":
				if (errz.Type().Kind() == reflect.String) {
					*errors = append(*errors, errz.StructField() + " must be more than or equal to " + errz.Param() + " character(s).")
				} else {
					*errors = append(*errors, errz.StructField() + " must be larger than " + errz.Param() + ".")
				}
			case "max":
				if (errz.Type().Kind() == reflect.String) {
					*errors = append(*errors, errz.StructField() + " must be lesser than or equal to " + errz.Param() + " character(s).")
				} else {
					*errors = append(*errors, errz.StructField() + " must be smaller than " + errz.Param() + ".")
				}
			default:
				*errors = append(*errors, errz.StructField() + " is invalid.")
		}		
	}

	return
}

// Merge two map string interface
func MergeMapString(mp1 map[string]interface{}, mp2 map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})
	for k, v := range mp1 {
        if _, ok := mp1[k]; ok {
            result[k] = v          
        }
    }

    for k, v := range mp2 {
        if _, ok := mp2[k]; ok {
            result[k] = v
        }
	}
	
	return result;
}