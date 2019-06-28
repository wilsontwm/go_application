package controllers

import (
	"os"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"html/template"
	"path/filepath"
	"net/url"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

var viewPath = "views"
var templates *template.Template
var restURL *url.URL
var appURL string
var appName string
var store *sessions.CookieStore
var cookieHashKey []byte
var cookieBlockKey []byte
var sCookie *securecookie.SecureCookie

func init() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Print("Error loading .env file", err)
	}

	templates, _ = GetTemplates()
	appName = os.Getenv("app_name")
	appURL = os.Getenv("app_url")
	restURL, _ = url.ParseRequestURI(appURL)
	store = sessions.NewCookieStore([]byte(os.Getenv("session_key")))
	//cookieHashKey = []byte(os.Getenv("cookie_hash_key"))
	//cookieBlockKey = []byte(os.Getenv("cookie_block_key"))
	sCookie = securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
}

func GetTemplates() (templates *template.Template, err error) {
	var allFiles []string
	
	// Loop through all the files in the views folder including subfolders
	err = filepath.Walk(viewPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		} 

		return nil
	})

	if err != nil {		
		log.Print("Error walking the file path", err)
	}

	templates, err = template.New("").ParseFiles(allFiles...)
	
	if err != nil {
		log.Print("Error parsing template files", err)
	}

    return
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string, cookieValue string) {
	value := cookieValue

	if encoded, err := sCookie.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name: cookieName,
			Value: encoded,			
			Path:  "/",
			// true means no scripts, http requests only
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
	}
}

func ClearCookieHandler(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie {
		Name: cookieName,
		Value: "",
		Path: "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string) (cookieValue string) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		if err = sCookie.Decode(cookieName, cookie.Value, &cookieValue); err == nil {
			return
		}
	}

	return 
}
