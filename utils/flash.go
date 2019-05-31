package utils

import (
  "net/http"
  "github.com/gorilla/sessions"
)

// Set the error/success flash message depends on the success state of the response
func SetErrorSuccessFlash(session *sessions.Session, w http.ResponseWriter, r *http.Request, resp map[string]interface{}) {
  // Set flash	
  var messages []interface{}
  
  if(resp["errors"] != nil) {
    messages = resp["errors"].([]interface{})
  } else {
    msg := resp["message"].(string)
    messages = append(messages, msg)
  }

  var tag string
  if(resp["success"].(bool)) {
		tag = "success"
	} else {	
    tag = "errors"
  }

  for _, message := range messages {
    session.AddFlash(message, tag)
  }

  session.Save(r, w)
}
