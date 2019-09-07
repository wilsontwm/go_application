package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"os"
	"github.com/joho/godotenv"
	"net/http"
	"app/api"
	"app/middleware"
)

func main() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	router := mux.NewRouter()

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
	apiProfileRoutes.HandleFunc("/edit", api.EditProfile).Methods("POST")
	apiProfileRoutes.HandleFunc("/edit/password", api.EditPassword).Methods("POST")
	apiProfileRoutes.HandleFunc("/upload/picture", api.UploadPicture).Methods("POST")
	apiProfileRoutes.HandleFunc("/delete/picture", api.DeletePicture).Methods("POST")

	// Invitation routes (incoming)
	apiInvitedRoutes := apiAuthenticatedRoutes.PathPrefix("/invite/incoming").Subrouter()
	apiInvitedRoutes.HandleFunc("", api.IndexInvitationFromCompany).Methods("GET")
	apiInvitedRoutes.HandleFunc("/{id}", api.ShowInvitationFromCompany).Methods("GET")
	apiInvitedRoutes.HandleFunc("/{id}/respond", api.RespondCompanyInvitationRequest).Methods("POST")
	
	// Company routes
	apiCompanyRoutes := apiAuthenticatedRoutes.PathPrefix("/company").Subrouter()
	apiCompanyRoutes.HandleFunc("", api.IndexCompany).Methods("GET")
	apiCompanyRoutes.HandleFunc("/store", api.CreateCompany).Methods("POST")
	apiCompanyRoutes.HandleFunc("/getUniqueSlug", api.GetUniqueSlug).Methods("GET")
	apiCompanyRoutes.HandleFunc("/{id}/show", api.ShowCompany).Methods("GET")
	apiCompanyRoutes.HandleFunc("/{id}/update", api.EditCompany).Methods("PATCH")
	apiCompanyRoutes.HandleFunc("/{id}/delete", api.DeleteCompany).Methods("DELETE")

	// Company invitation request routes (outgoing)
	apiCompanyRoutes.HandleFunc("/{id}/invite", api.InviteToCompany).Methods("POST")
	apiCompanyRoutes.HandleFunc("/{id}/invite/list", api.IndexInviteToCompany).Methods("GET")
	apiCompanyRoutes.HandleFunc("/{id}/invite/{invitationID}", api.ShowCompanyInvitationRequest).Methods("GET")
	apiCompanyRoutes.HandleFunc("/{id}/invite/{invitationID}/delete", api.DeleteCompanyInvitationRequest).Methods("DELETE")
	
	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	log.Println("Server started and running at port", port)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS(headers, methods, origins)(router)))
}