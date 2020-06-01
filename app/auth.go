package app

import (
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"my-contacts/models"
	utl "my-contacts/utils"
	"net/http"
	"os"
	"strings"
)

/*
Structure of a basic middleware

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff
		next.ServeHTTP(w, r)
	})
}
*/

// JwtAuthentication middleware
var JwtAuthentication = func (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Slice of paths that does not require authentication
		notAuth := []string{"/api/v1/user/new", "/api/v1/user/login"}
		requestPath := r.URL.Path // current path

		// check if request does not need authentication
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string] interface{})
		tokenHeader := r.Header.Get("Authorization") // retrieve the token

		if tokenHeader == "" { // check if token is missing
			response = utl.Message(false, "auth token missing")
			w.Header().Add("Content-Type", "application/json") // Headers should come before Status
			w.WriteHeader(http.StatusForbidden)
			// fmt.Println("check here")
			utl.Respond(w, response)
			return
		}

		// Fetching the token
		splitText := strings.Split(tokenHeader, " ")
		if len(splitText) != 2 {
			response = utl.Message(false, "Invalid/Malformed auth token")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			utl.Respond(w, response)
			return
		}

		tokenValue := splitText[1] // Get the token
		tk := &models.Token{}

		// Decode our token
		token, err := jwt.ParseWithClaims(tokenValue, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// Malformed token
		if err != nil {
			response = utl.Message(false, "Malformed authentication token")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			utl.Respond(w, response)
			return
		}

		// Invalid token
		if !token.Valid {
			response = utl.Message(false, "Token is invalid")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			utl.Respond(w, response)
			return
		}

		// Everything went well
		// fmt.Sprintf("User %s", string(tk.UserId))
		ctx := context.WithValue(r.Context(), "user", tk.UserId) // ctx variable to store claims
		r = r.WithContext(ctx) // Setting context to request
		next.ServeHTTP(w, r)
	})
}
