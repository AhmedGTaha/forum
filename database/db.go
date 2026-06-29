package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Blank import to register the driver
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

// DB is global so our HTTP handlers can access it later
var DB *sql.DB

func InitDB() {
	var err error
	// Opens the database file, creating it if it doesn't exist
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// SQL query to create all necessary tables
	createTablesQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		content TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(post_id) REFERENCES posts(id)
	);
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		comment_id INTEGER,
		is_like BOOLEAN,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(comment_id) REFERENCES comments(id)
	);
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER,
		expires_at DATETIME,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	// executes the query
	_, err = DB.Exec(createTablesQuery)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	log.Println("Database and tables initialized successfully!")
}

// searches for one user using their email
func GetUserByEmail(email string) (User, error) {
	var user User
	
	query := `
		SELECT id, username, email, password
		FROM users
		WHERE email = ?
	`

	// 1. Runs the SELECT query
	// 2. Sends the email into the ? placeholder
	// 3. Copies the result into the user struct
	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}