package handlers

import (
	"net/http"

	"forum/database"
)

// GetLoggedInUser checks whether the request has a valid session cookie
func GetLoggedInUser(r *http.Request) (database.User, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return database.User{}, false
	}

	if cookie.Value == "" {
		return database.User{}, false
	}

	user, err := database.GetUserBySessionID(cookie.Value)
	if err != nil {
		return database.User{}, false
	}

	return user, true
}