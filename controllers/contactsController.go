package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"my-contacts/models"
	utl "my-contacts/utils"
	"net/http"
	"strconv"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) // Get the user that send the request
	contactPointer := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contactPointer)
	if err != nil {
		utl.Respond(w, utl.Message(false, "Error while decoding request body"))
		return
	}

	contactPointer.UserId = user
	resp := contactPointer.Create()
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
