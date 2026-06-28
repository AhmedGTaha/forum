package database

import (
	"database/sql"
	"log"

	// This blank import is required to register the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// DB is a global variable so other packages (like handlers) can access the database connection
var DB *sql.DB

func InitDB() {
	var err error
	// Open the database (creates forum.db in your root folder if it's not there)
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create the users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Database initialized and users table created successfully!")
}