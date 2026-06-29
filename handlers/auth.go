package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/database"

	"golang.org/x/crypto/bcrypt"
)

// LoginPageHandler displays the login page.
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/login.html")
	if err != nil {
		http.Error(w, "Could not load login page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not render login page", http.StatusInternalServerError)
		return
	}
}

// LoginHandler handles the submitted login form.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || strings.TrimSpace(password) == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// LoginDispatcher chooses the correct login handler based on the HTTP method.
func LoginDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		LoginPageHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		LoginHandler(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// RegisterPageHandler displays the registration page.
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/register.html")
	if err != nil {
		http.Error(w, "Could not load registration page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not render registration page", http.StatusInternalServerError)
		return
	}
}

// RegisterHandler handles the submitted registration form.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if username == "" || email == "" || strings.TrimSpace(password) == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not secure password", http.StatusInternalServerError)
		return
	}

	queryInsertUser := `
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)
	`

	_, err = database.DB.Exec(queryInsertUser, username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// RegisterDispatcher chooses the correct registration handler based on the HTTP method.
func RegisterDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RegisterPageHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		RegisterHandler(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}