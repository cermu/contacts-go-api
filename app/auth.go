package app

import (
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"my-contacts/models"
	utl "my-contacts/utils"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/new", "/api/user/login"} // Slice of paths that does not require authentication
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
			w.WriteHeader(http.StatusForbidden)
			// w.Header().Add("Content-Type", "application/json")
			utl.Respond(w, response)
			return
		}

		// Fetching the token
		splitText := strings.Split(tokenHeader, " ")
		if len(splitText) != 2 {
			response = utl.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			utl.Respond(w, response)
			return
		}

		tokenValue := splitText[1] // Get the token
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenValue, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// Malformed token
		if err != nil {
			response = utl.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			utl.Respond(w, response)
			return
		}

		// Invalid token
		if !token.Valid {
			response = utl.Message(false, "Token is invalid")
			w.WriteHeader(http.StatusForbidden)
			utl.Respond(w, response)
			return
		}

		// Everything went well
		fmt.Sprintf("User %s", string(tk.UserId))
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
