package handlers

import (
	"fmt"
	"forum/database"
	"net/http"
	"html/template"
	"golang.org/x/crypto/bcrypt"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the template file
    tmpl, err := template.ParseFiles("ui/register.html")
    if err != nil {
        http.Error(w, "Could not load template", http.StatusInternalServerError)
        return
    }

    // Execute the template and write to the response
    tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// r.FormValue: extract data from an HTML form submission
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // create a secure hash
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Insert new user into database
	queryInsertUser := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err = database.DB.Exec(queryInsertUser, username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User registered successfully!")
}

func RegisterDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RegisterPageHandler(w, r) // here it displays the page
	} else if r.Method == http.MethodPost {
		// it sends data from the html form
		RegisterHandler(w, r) // here it sends info to db
	} else {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}