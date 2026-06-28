package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"log"
	"net/http" // for the web server functionality
)

func main () {
	// Initialize the database
	err := database.InitDB()
	fmt.Println("Database initialized successfully!")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// run the function whenever the url is /
	http.HandleFunc("/", handlers.HomeHandler)

	// This starts the server
	fmt.Println("Server is starting on http://localhost:8080...")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}