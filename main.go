package main

import (
	"github.com/gorilla/mux"
	"log"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"app/controllers"
	"app/api"
	"app/middleware"
)

func main() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	router := mux.NewRouter()
	// Routes
	// Pages routes
	router.HandleFunc("/", controllers.HelloPage).Methods("GET")
	// Login / register routes	
	nonAuthenticatedRoutes := router.PathPrefix("").Subrouter()
	nonAuthenticatedRoutes.Use(middleware.LogTime(), middleware.Second())
	nonAuthenticatedRoutes.HandleFunc("/login", controllers.LoginPage).Methods("GET")
	nonAuthenticatedRoutes.HandleFunc("/login", controllers.LoginSubmit).Methods("POST")
	nonAuthenticatedRoutes.HandleFunc("/signup", controllers.SignupPage).Methods("GET")
	nonAuthenticatedRoutes.HandleFunc("/signup", controllers.SignupSubmit).Methods("POST")
	nonAuthenticatedRoutes.HandleFunc("/resendactivation", controllers.ResendActivationPage).Methods("GET")	
	nonAuthenticatedRoutes.HandleFunc("/resendactivation", controllers.ResendActivationSubmit).Methods("POST")
	nonAuthenticatedRoutes.HandleFunc("/activate/{code}", controllers.ActivateAccountPage).Methods("GET")	
	nonAuthenticatedRoutes.HandleFunc("/forgetpassword", controllers.ForgetPasswordPage).Methods("GET")	
	nonAuthenticatedRoutes.HandleFunc("/forgetpassword", controllers.ForgetPasswordSubmit).Methods("POST")
	nonAuthenticatedRoutes.HandleFunc("/resetpassword/{code}", controllers.ResetPasswordPage).Methods("GET")	
	nonAuthenticatedRoutes.HandleFunc("/resetpassword/{code}", controllers.ResetPasswordSubmit).Methods("POST")
	
	// REST routes
	apiRoutes := router.PathPrefix("/api").Subrouter()
	apiRoutes.HandleFunc("/login", api.Login).Methods("POST")
	apiRoutes.HandleFunc("/signup", api.Signup).Methods("POST")
	apiRoutes.HandleFunc("/resendactivation", api.ResendActivation).Methods("POST")
	apiRoutes.HandleFunc("/activateaccount", api.ActivateAccount).Methods("POST")
	apiRoutes.HandleFunc("/forgetpassword", api.ForgetPassword).Methods("POST")
	apiRoutes.HandleFunc("/resetpassword", api.ResetPassword).Methods("POST")

	// Asset files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	// Custom 404 page
	router.NotFoundHandler = http.HandlerFunc(controllers.Custom404Page)

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	log.Println("Server started and running at port", port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}
