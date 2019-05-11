package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"app/controllers"
	"app/api"
)

func main() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	router := mux.NewRouter()
	// Routes
	// Pages routes
	router.HandleFunc("/", controllers.HelloPage).Methods("GET")
	// Authenticate routes
	router.HandleFunc("/login", controllers.LoginPage).Methods("GET")
	router.HandleFunc("/signup", controllers.SignupPage).Methods("GET")
	router.HandleFunc("/signup", controllers.SignupSubmit).Methods("POST")
    
	// REST routes
	router.HandleFunc("/api/signup", api.Signup).Methods("POST")

	// Asset files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Server started and running at port", port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}
