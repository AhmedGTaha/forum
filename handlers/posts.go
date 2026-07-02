package handlers

import (
	"html/template"
	"net/http"
)

// This is the data we send to the HTML page
type CreatePostPageData struct {
	Username string
}

// Only logged-in users are allowed to access this page.
func CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the user is logged in
	user, loggedIn := GetLoggedInUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := CreatePostPageData{
		Username: user.Username,
	}

	tmpl, err := template.ParseFiles("ui/create-post.html")
	if err != nil {
		http.Error(w, "Could not load create post page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not render create post page", http.StatusInternalServerError)
		return
	}
}

// CreatePostDispatcher chooses the correct handler based on the HTTP method.
func CreatePostDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CreatePostPageHandler(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
