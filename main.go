package main

import (
	"fmt"
	"log"
	"net/http"
	
	"forum/database" // Make sure this matches your module name in go.mod
	"forum/handlers"
)

func main() {
	// Initialize the database before starting the server
	database.InitDB()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterDispatcher)

	fmt.Println("Server is starting on http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}