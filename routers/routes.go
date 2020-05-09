package routers

import (
	"github.com/gorilla/mux"
	"my-contacts/app"
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

// NewRouter function for registering endpoints and their handlers.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(app.JwtAuthentication) // Attach the JWT middleware
	api := router.PathPrefix("/api/v1").Subrouter()

	for _, route := range routes {
		api.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

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
