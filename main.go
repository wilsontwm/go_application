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
	nonAuthenticatedRoutes := router.PathPrefix("").Subrouter()
	
	// Pages routes
	nonAuthenticatedRoutes.HandleFunc("/", controllers.WelcomePage).Methods("GET").Name("welcome")
	nonAuthenticatedRoutes.HandleFunc("/noaccess", controllers.Custom403Page).Name("error_403")

	// Login / register routes	
	nonAuthenticatedRoutes.HandleFunc("/login", controllers.LoginPage).Methods("GET").Name("login")
	nonAuthenticatedRoutes.HandleFunc("/login", controllers.LoginSubmit).Methods("POST").Name("login_submit")
	nonAuthenticatedRoutes.HandleFunc("/logout", controllers.LogoutSubmit).Methods("GET").Name("logout")
	nonAuthenticatedRoutes.HandleFunc("/signup", controllers.SignupPage).Methods("GET").Name("signup")
	nonAuthenticatedRoutes.HandleFunc("/signup", controllers.SignupSubmit).Methods("POST").Name("signup_submit")
	nonAuthenticatedRoutes.HandleFunc("/resendactivation", controllers.ResendActivationPage).Methods("GET").Name("resend_activation")	
	nonAuthenticatedRoutes.HandleFunc("/resendactivation", controllers.ResendActivationSubmit).Methods("POST").Name("resend_activation_submit")
	nonAuthenticatedRoutes.HandleFunc("/activate/{code}", controllers.ActivateAccountPage).Methods("GET").Name("activate_account")	
	nonAuthenticatedRoutes.HandleFunc("/forgetpassword", controllers.ForgetPasswordPage).Methods("GET").Name("forget_password")	
	nonAuthenticatedRoutes.HandleFunc("/forgetpassword", controllers.ForgetPasswordSubmit).Methods("POST").Name("forget_password_submit")
	nonAuthenticatedRoutes.HandleFunc("/resetpassword/{code}", controllers.ResetPasswordPage).Methods("GET").Name("reset_password")	
	nonAuthenticatedRoutes.HandleFunc("/resetpassword/{code}", controllers.ResetPasswordSubmit).Methods("POST").Name("reset_password_submit")
	
	authenticatedRoutes := router.PathPrefix("/dashboard").Subrouter()
	authenticatedRoutes.Use(middleware.CheckAuth())
	authenticatedRoutes.HandleFunc("", controllers.DashboardPage).Methods("GET").Name("dashboard")
	authenticatedRoutes.HandleFunc("/profile/edit", controllers.EditProfilePage).Methods("GET").Name("profile_edit")
	
	// REST routes
	apiRoutes := router.PathPrefix("/api").Subrouter()
	apiRoutes.HandleFunc("/login", api.Login).Methods("POST")
	apiRoutes.HandleFunc("/signup", api.Signup).Methods("POST")
	apiRoutes.HandleFunc("/resendactivation", api.ResendActivation).Methods("POST")
	apiRoutes.HandleFunc("/activateaccount", api.ActivateAccount).Methods("POST")
	apiRoutes.HandleFunc("/forgetpassword", api.ForgetPassword).Methods("POST")
	apiRoutes.HandleFunc("/resetpassword", api.ResetPassword).Methods("POST")

	apiAuthenticatedRoutes := apiRoutes.PathPrefix("/dashboard").Subrouter()
	apiAuthenticatedRoutes.Use(middleware.JwtAuthentication())

	// Profiles routes
	apiProfileRoutes := apiAuthenticatedRoutes.PathPrefix("/profile").Subrouter()
	apiProfileRoutes.HandleFunc("/get", api.GetProfile).Methods("GET")

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
