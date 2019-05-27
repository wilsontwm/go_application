package controllers

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"html/template"
	"path/filepath"
	"net/url"
)

var viewPath = "views"
var templates *template.Template
var restURL *url.URL
var appURL string
var appName string

func init() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Print("Error loading .env file", err)
	}

	templates, _ = GetTemplates()
	appName = os.Getenv("app_name")
	appURL = os.Getenv("app_url")
	restURL, _ = url.ParseRequestURI(appURL)
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
