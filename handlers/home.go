package handlers

import (
	"net/http"
	"text/template" // for the parsing thing
)

type HomePageData struct {
	IsLoggedIn bool
	Username   string
}

// w: to write our response back to the user's browser
// r: incoming request from the user
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

    data := HomePageData{}

    user, loggedIn := GetLoggedInUser(r)
    if loggedIn {
        data.IsLoggedIn = true
        data.Username = user.Username
    }

	tmpl, err := template.ParseFiles("ui/index.html")
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not render home page", http.StatusInternalServerError)
		return
	}
}