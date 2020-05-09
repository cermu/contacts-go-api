package main

import (
	"fmt"
	"log"
	rte "my-contacts/routers"
	utl "my-contacts/utils"
	"net/http"
	"os"
)

func main() {
	utl.WriteToFile("*********************" +
		"****************")
	utl.WriteToFile("INFO | Starting the application...")

	router := rte.NewRouter()

	port := os.Getenv("application_port")
	if port == "" {
		port = "8000"
	}

	utl.WriteToFile(fmt.Sprintf("INFO | Application is running on port: %v", port))

	err := http.ListenAndServe(":" + port, router) // Launch the app
	if err != nil {
		// fmt.Print(err)
		log.Fatal(err)
	}
}
