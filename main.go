package main

import (
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Server started and running at port", port)

	err := http.ListenAndServe(":" + port, router)
	if err != nil {
		log.Print("Server running error: ", err) 
	}
}