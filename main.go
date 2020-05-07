package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"my-contacts/app"
	"my-contacts/controllers"
	"net/http"
	"os"
)

func main() {
	log.Println("*********************************")
	log.Println("Starting the application...")
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

	log.Printf("Application is running on port: %v \n", port)

	err := http.ListenAndServe(":" + port, router) // Launch the app
	if err != nil {
		fmt.Print(err)
	}
}
