package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"my-contacts/models"
	utl "my-contacts/utils"
	"net/http"
	"strconv"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) // Get the user that sent the request
	contactPointer := &models.Contact{}

	requestBody1, requestBody2 := utl.RetrieveRequestBody(r)
	r.Body = requestBody2

	err := json.NewDecoder(r.Body).Decode(contactPointer)
	if err != nil {
		utl.Respond(w, utl.Message(false, "Error while decoding request body"))
		return
	}

	utl.WriteToFile(fmt.Sprintf("INFO | Type = Request | " +
		"Source = %v | Target System = | Request Body = %v", r.Host, requestBody1))

	contactPointer.UserId = user
	resp := contactPointer.Create()
	responseBody, _ := json.Marshal(resp)

	utl.WriteToFile(fmt.Sprintf("INFO | Type = Response | Target System = | Response Body = %s", responseBody))
	utl.Respond(w, resp)
}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// fmt.Println(params)
	id, err := strconv.Atoi(params["userId"])
	if err != nil {
		utl.Respond(w, utl.Message(false, "There was an error in your request"))
		return
	}

	data := models.GetContacts(uint(id))
	resp := utl.Message(true, "success")
	resp["data"] = data
	utl.Respond(w, resp)
}
