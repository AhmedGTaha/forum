package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

// post we want to display
type Post struct {
	ID        int
	UserID    int
	Username  string
	Title     string
	Content   string
	CreatedAt string
}

var DB *sql.DB

func InitDB() error {
	var err error

	DB, err = sql.Open("sqlite3", "file:forum.db?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	DB.SetMaxOpenConns(1)

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	err = createTables()
	if err != nil {
		return err
	}

	fmt.Println("Database initialized successfully!")
	return nil
}

func createTables() error {
	queries := []string{
		createUsersTable,
		createSessionsTable,
		createCategoriesTable,
		createPostsTable,
		createPostCategoriesTable,
		createCommentsTable,
		createPostReactionsTable,
		createCommentReactionsTable,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}

	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

const createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL UNIQUE,
	expires_at TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`

const createCategoriesTable = `
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);
`

const createPostsTable = `
CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`

const createPostCategoriesTable = `
CREATE TABLE IF NOT EXISTS post_categories (
	post_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	PRIMARY KEY (post_id, category_id),
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
`

const createCommentsTable = `
CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
`

const createPostReactionsTable = `
CREATE TABLE IF NOT EXISTS post_reactions (
	user_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
	value INTEGER NOT NULL CHECK (value IN (1, -1)),
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, post_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
`

const createCommentReactionsTable = `
CREATE TABLE IF NOT EXISTS comment_reactions (
	user_id INTEGER NOT NULL,
	comment_id INTEGER NOT NULL,
	value INTEGER NOT NULL CHECK (value IN (1, -1)),
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, comment_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);
`

func GetUserByEmail(email string) (User, error) {
	var user User

	query := `
		SELECT id, username, email, password
		FROM users
		WHERE email = ?
	`

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

func CreateSession(userID int) (string, time.Time, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	expiresAtText := expiresAt.Format(time.RFC3339Nano)

	tx, err := DB.Begin()
	if err != nil {
		return "", time.Time{}, err
	}

	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		return "", time.Time{}, err
	}

	query := `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)
	`

	_, err = tx.Exec(query, sessionID, userID, expiresAtText)
	if err != nil {
		tx.Rollback()
		return "", time.Time{}, err
	}

	err = tx.Commit()
	if err != nil {
		return "", time.Time{}, err
	}

	return sessionID, expiresAt, nil
}

func GetUserBySessionID(sessionID string) (User, error) {
	var user User

	query := `
		SELECT users.id, users.username, users.email, users.password
		FROM sessions
		INNER JOIN users ON users.id = sessions.user_id
		WHERE sessions.id = ?
		AND sessions.expires_at > ?
	`

	now := time.Now().UTC().Format(time.RFC3339Nano)

	err := DB.QueryRow(query, sessionID, now).Scan(
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

func DeleteSession(sessionID string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteExpiredSessions() error {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err := DB.Exec("DELETE FROM sessions WHERE expires_at <= ?", now)
	if err != nil {
		return err
	}

	return nil
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// Receives the logged-in user ID, post title, and post content
func CreatePost(userID int, title string, content string) (int64, error) {
	query := `
		INSERT INTO posts (user_id, title, content)
		VALUES (?, ?, ?)
	`

	result, err := DB.Exec(query, userID, title, content)
	if err != nil {
		return 0, err
	}

	// This gets the ID of the post that was just created
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func GetAllPosts() ([]Post, error) {
	// This joins posts with their authors, Without this, we would only have: user_id = 1 but now we have: Posted by omar
	query := `
		SELECT posts.id, posts.user_id, users.username, posts.title, posts.content, posts.created_at
		FROM posts
		INNER JOIN users ON posts.user_id = users.id
		ORDER BY posts.created_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	// makes sure the database rows are closed after we finish reading them
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// This checks if something went wrong while looping through rows
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}