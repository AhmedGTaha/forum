package main

import (
	"fmt"           // format text messages
	"html/template" // load and render HTML templates
	"log"           // print server messages
	"net/http"      // tools to build the server
)

// HomePageData is the data we send to home.html
type HomePageData struct {
	Title   string // Page title
	Message string // Welcome message or status info
}

// homeHandler runs when the browser visits the root path "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user visited exactly "/" or something else
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Prepare data for the template
	pageData := HomePageData{
		Title:   "Forum",
		Message: "Welcome to the Forum!",
	}

	// Parse the home.html template from the templates folder
	tmpl, err := template.ParseFiles("web/templates/home.html")
	if err != nil {
		log.Println("Template parsing error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the HTML page to the browser with pageData
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// pingHandler is a simple health check endpoint
// It responds with "pong" to confirm the server is running
func pingHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests to this endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}

func main() {
	// mux = multiplexer (router)
	mux := http.NewServeMux()

	// Serve static files (CSS, JS, images) from the web/static folder
	// Any URL starting with "/static/" will serve files from web/static/
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Register route handlers
	mux.HandleFunc("/ping", pingHandler)   
	mux.HandleFunc("/", homeHandler)

	port := ":8080"
	log.Printf("Server starting at http://localhost%s", port)

	// Start the server
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}