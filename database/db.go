package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error

	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	// SQL query to create a table for users
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password TEXT
	);`

	_, err = DB.Exec(createUsersTable)
	return err
}
