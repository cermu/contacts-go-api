package routers

import (
	"my-contacts/controllers"
	"net/http"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		Name:        "Authenticate",
		Method:      "POST",
		Pattern:     "/user/login",
		HandlerFunc: controllers.Authenticate,
	},
	Route{
		Name:        "CreateAccount",
		Method:      "POST",
		Pattern:     "/user/new",
		HandlerFunc: controllers.CreateAccount,
	},
	Route{
		Name:        "CreateContact",
		Method:      "POST",
		Pattern:     "/contact/new",
		HandlerFunc: controllers.CreateContact,
	},
	Route{
		Name:        "GetContactsFor",
		Method:      "GET",
		Pattern:     "/user/{userId}/contacts",
		HandlerFunc: controllers.GetContactsFor,
	},
}
