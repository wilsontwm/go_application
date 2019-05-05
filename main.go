package main

import (
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"path/filepath"
)

var viewPath = "views"
var templates *template.Template

func getTemplates() (templates *template.Template, err error) {
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

func init() {
    templates, _ = getTemplates()
}

func main() {
	router := mux.NewRouter()
	// Routes
	router.HandleFunc("/", hello).Methods("GET")

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Server started and running at port", port)

	log.Fatal(http.ListenAndServe(":" + port, router))
}

func hello(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Homepage",
	}

	err := templates.ExecuteTemplate(w, "home_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
