package handlers

import (
	"net/http"
	"text/template" // for the parsing thing
)

// w: to write our response back to the user's browser
// r: incoming request from the user
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the template file
    tmpl, err := template.ParseFiles("ui/index.html")
    if err != nil {
        http.Error(w, "Could not load template", http.StatusInternalServerError)
        return
    }

    // Execute the template and write to the response
    tmpl.Execute(w, nil)
}