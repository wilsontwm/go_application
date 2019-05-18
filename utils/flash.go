package utils

import (
  "encoding/base64"
  "net/http"
  "time"  
  "encoding/json"
)

// Set the error/success flash message depends on the success state of the response
func SetErrorSuccessFlash(w http.ResponseWriter, resp map[string]interface{}) {
	if(resp["success"].(bool)) {
		successFlash := []byte(resp["message"].(string))
		SetFlash(w, "success", successFlash)
	} else {
		// Set flash		
		errors := resp["errors"].([] interface{})
		errorJson, _ := json.Marshal(errors)
		errorFlash := []byte(errorJson)
		SetFlash(w, "errors", errorFlash)
	}
}

// Set the flash message into cookie
func SetFlash(w http.ResponseWriter, name string, value []byte) {
  c := &http.Cookie{Name: name, Value: encode(value)}
  http.SetCookie(w, c)
}

// Get the flash message from cookie
func GetFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
  c, err := r.Cookie(name)
  if err != nil {
    switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
    }
  }
  value, err := decode(c.Value)
  if err != nil {
    return nil, err
  }
  dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
  http.SetCookie(w, dc)
  return value, nil
}

func encode(src []byte) string {
  return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
  return base64.URLEncoding.DecodeString(src)
}