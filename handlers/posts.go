package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/database"
)

// This is the data we send to the HTML page.
type CreatePostPageData struct {
	Username string
}

// Only logged-in users are allowed to access this page.
func CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in.
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

// CreatePostHandler handles the submitted create post form.
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// This handler should only accept POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in.
	user, loggedIn := GetLoggedInUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Read and clean the form values.
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	// Validate that the user did not submit empty data.
	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Save the post in the database using the logged-in user's ID.
	_, err := database.CreatePost(user.ID, title, content)
	if err != nil {
		http.Error(w, "Could not create post", http.StatusInternalServerError)
		return
	}

	// Redirect after successful form submission.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreatePostDispatcher chooses the correct handler based on the HTTP method.
func CreatePostDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CreatePostPageHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		CreatePostHandler(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}