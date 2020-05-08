package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"my-contacts/app"
	"my-contacts/controllers"
	utl "my-contacts/utils"
	"net/http"
	"os"
)

func main() {
	utl.WriteToFile("*************************************")
	utl.WriteToFile("INFO | Starting the application...")
	router := mux.NewRouter().StrictSlash(true)
	router.Use(app.JwtAuthentication) // Attach the JWT auth middleware
	api := router.PathPrefix("/api/v1").Subrouter()

	// Registering endpoints and their corresponding request handlers
	api.HandleFunc("/user/new", controllers.CreateAccount).Methods("POST")
	api.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")
	api.HandleFunc("/user/{userId}/contacts", controllers.GetContactsFor).Methods("GET")
	api.HandleFunc("/contacts/new", controllers.CreateContact).Methods("POST")

	port := os.Getenv("application_port")
	if port == "" {
		port = "8000"
	}

	utl.WriteToFile(fmt.Sprintf("INFO | Application is running on port: %v", port))

	err := http.ListenAndServe(":" + port, router) // Launch the app
	if err != nil {
		fmt.Print(err)
	}
}
