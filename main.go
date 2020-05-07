package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"my-contacts/app"
	"my-contacts/controllers"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	currentDate := time.Now()
	y, m, d := currentDate.Date()
	truncateMonth := m.String()[:3]
	if err := ensureDir("./logs/" + truncateMonth); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logFile := strconv.Itoa(d) + strconv.Itoa(int(m)) + strconv.Itoa(y)

	f, er := os.OpenFile("./logs/" + truncateMonth + "/api_" + logFile +".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if er != nil{
		log.Fatalf("Errpr openning file %v", er)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	log.Println("*********************************")
	log.Println("INFO | Starting the application...")
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

	log.Printf("INFO | Application is running on port: %v \n", port)

	err := http.ListenAndServe(":" + port, router) // Launch the app
	if err != nil {
		fmt.Print(err)
	}
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
