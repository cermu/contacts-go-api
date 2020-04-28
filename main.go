package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"my-contacts/app"
	"my-contacts/controllers"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting the application...")
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // Attach the JWT auth middleware

	// Registering endpoints and their corresponding request handlers
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	port := os.Getenv("application_port")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Application is running on port: %v \n", port)

	err := http.ListenAndServe(":" + port, router) // Launch the app
	if err != nil {
		fmt.Print(err)
	}
}
