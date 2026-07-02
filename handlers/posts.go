package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/database"
)

// This is the data we send to the HTML page.
type CreatePostPageData struct {
	Username   string
	Categories []database.Category
}

// Only logged-in users are allowed to access this page.
func CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	user, loggedIn := GetLoggedInUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	categories, err := database.GetAllCategories()
	if err != nil {
		http.Error(w, "Could not load categories", http.StatusInternalServerError)
		return
	}

	data := CreatePostPageData{
		Username:   user.Username,
		Categories: categories,
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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
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

	// Because the user can select multiple categories we need the full list
	categoryValues := r.Form["categories"]
	if len(categoryValues) == 0 {
		http.Error(w, "At least one category must be selected", http.StatusBadRequest)
		return
	}

	var categoryIDs []int

	for _, categoryValue := range categoryValues {
		// Convert the category ID from string to int
		categoryID, err := strconv.Atoi(categoryValue)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	_, err = database.CreatePost(user.ID, title, content, categoryIDs)
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