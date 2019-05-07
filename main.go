package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"app/controllers"
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
	router.HandleFunc("/login", controllers.LoginPage).Methods("GET")
	router.HandleFunc("/signup", controllers.SignupPage).Methods("GET")

	// REST routes
	
	// Asset files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Server started and running at port", port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}
