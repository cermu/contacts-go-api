package controllers

import (
	"encoding/json"
	"my-contacts/models"
	utl "my-contacts/utils"
	"net/http"
)

// CreateAccount handler for creating new users
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	accountPointer := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(accountPointer) // decode the request body into a struct
	if err != nil {
		utl.Respond(w, utl.Message(false, "Invalid request"))
		return
	}

	resp := accountPointer.Create() // Create an account
	utl.Respond(w, resp)
}

// Authenticate handler for authenticating users
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	accountPointer := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(accountPointer) // decode the request body into a struct
	if err != nil {
		utl.Respond(w, utl.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(accountPointer.Email, accountPointer.Password)
	utl.Respond(w, resp)
}

// UserLogout handler for logging out users
var UserLogout = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	resp := models.Logout(userId)

	utl.Respond(w, resp)
}
