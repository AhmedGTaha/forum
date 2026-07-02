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

// LoginHandler handles the submitted login form
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

	// delete old session for the same user
	// create a new random session ID
	// store it in the sessions table
	// return the session ID and expiration time
	sessionID, expiresAt, err := database.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	// This prepares the browser cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	// This sends the cookie to the browser in the HTTP response
	http.SetCookie(w, &cookie)

	// After login, the user goes back to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

// LogoutHandler logs the user out by deleting the session and clearing the cookie.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		err = database.DeleteSession(cookie.Value)
		if err != nil {
			http.Error(w, "Could not delete session", http.StatusInternalServerError)
			return
		}
	}

	expiredCookie := http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &expiredCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
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