package routers

import (
	"github.com/gorilla/mux"
	"my-contacts/app"
)

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
