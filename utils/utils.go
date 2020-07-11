package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Message function builds json messages
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond function responds with json message
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

// WriteToFile function creates the logs dir and a monthly sub dir in it.
// Then a log file is created to store the application logs.
func WriteToFile(msg string) {
	e := godotenv.Load() // Load .env file
	if e != nil {
		// fmt.Print(e)
		log.Fatalf("The following error occurred while loading env file: %s", e)
	}

	DEBUG, boolErr := strconv.ParseBool(os.Getenv("DEBUG"))
	if boolErr != nil {
		log.Fatalf("Error parseing bool: %v", boolErr)
	}

	if DEBUG {
		fmt.Println(msg)
	} else {
		currentDate := time.Now()
		y, m, d := currentDate.Date()
		truncateMonth := m.String()[:3]
		dirName := "./logs/" + truncateMonth
		logFile := strconv.Itoa(d) + strconv.Itoa(int(m)) + strconv.Itoa(y)

		// Creating the log directory
		dirErr := os.MkdirAll(dirName, os.ModePerm)
		if dirErr != nil || os.IsExist(dirErr) {
			log.Fatalf("Error creating directory: %v", dirErr)
		}

		// Creating the log file
		f, er := os.OpenFile(dirName+"/api_log"+logFile+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if er != nil {
			log.Fatalf("Error opening file: %v", er)
		}
		wrt := io.MultiWriter(os.Stdout, f)
		log.SetOutput(wrt)
		log.Println(msg)

		defer f.Close()
	}
}

// RetrieveRequestBody function that retrieves the body from POST request
func RetrieveRequestBody(r *http.Request) (string, io.ReadCloser) {
	// r.Body is a buffer, which means,
	// once it has been read, it cannot be read again.
	// Catch the body and restore it for other uses.
	// requestBody1 will be used for logging .
	// requestBody2 will be used for normal request processing.

	body, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		// fmt.Println(bodyErr)
		WriteToFile(fmt.Sprintf("ERROR | The following ioutil.ReadAll error occurred: %s", bodyErr))
	}
	requestBody1 := string(body)
	requestBody2 := ioutil.NopCloser(bytes.NewBuffer(body))
	return requestBody1, requestBody2
}
