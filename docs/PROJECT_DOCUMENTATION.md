# Forum Project Documentation

## 1. Project overview

This project is a small Go web application that aims to become a basic forum with:

- user registration
- login and sessions
- creating posts
- adding comments
- liking/disliking posts and comments
- filtering posts by category or ownership
- a simple SQLite database backend

The current codebase is still in the early stage. The app can start a web server, initialize a database, and register users.

---

## 2. Project goal

The original project requirements are described in [docs/README(1).md](README(1).md). The main objective is to build a forum that supports user interaction and persistence through SQLite.

The project is meant to be a learning project, so the implementation is intentionally simple and uses:

- Go's built-in net/http package
- HTML templates in the ui folder
- SQLite through the github.com/mattn/go-sqlite3 driver
- bcrypt for password hashing

---

## 3. Current project structure

```text
forum/
├── database/
│   └── db.go
├── docs/
│   ├── PROJECT_DOCUMENTATION.md
│   ├── README(1).md
│   └── fix.txt
├── handlers/
│   ├── auth.go
│   └── home.go
├── internal/
│   └── database/
│       └── db.go
├── ui/
│   ├── index.html
│   └── register.html
├── forum.db
├── go.mod
├── go.sum
├── main.go
```

### Important folders and files

- [main.go](../main.go)
  - application entry point
  - initializes the database
  - registers routes
  - starts the HTTP server

- [handlers/auth.go](../handlers/auth.go)
  - registration-related handlers
  - renders registration page
  - handles POST requests for new users
  - dispatches GET and POST requests for /register

- [handlers/home.go](../handlers/home.go)
  - serves the home page
  - renders the index template

- [database/db.go](../database/db.go)
  - the database package currently used by the app
  - creates the SQLite tables
  - exposes the global DB connection

- [internal/database/db.go](../internal/database/db.go)
  - duplicate database package implementation
  - currently not used by the app entry point
  - can be ignored for now or removed later to avoid confusion

- [ui/index.html](../ui/index.html)
  - homepage template

- [ui/register.html](../ui/register.html)
  - registration form template

- [docs/README(1).md](README(1).md)
  - original assignment/specification

- [docs/fix.txt](fix.txt)
  - Windows setup instructions for Go + SQLite + CGO

---

## 4. Runtime flow

### Startup sequence

1. [main.go](../main.go) runs.
2. It calls database.InitDB().
3. It registers the handlers:
   - / -> HomeHandler
   - /register -> RegisterDispatcher
4. It starts the server on port 8080.

### Request flow

- Visiting / shows the home page via HomeHandler.
- Visiting /register with GET shows the registration form.
- Submitting the form to /register with POST sends the data to RegisterHandler.
- RegisterHandler hashes the password and inserts the user into the users table.

---

## 5. Functions and where they live

| Function | Location | Responsibility | Current state |
|---|---|---|---|
| main | [main.go](../main.go) | Starts the server and wires routes | Implemented |
| RegisterPageHandler | [handlers/auth.go](../handlers/auth.go) | Renders the registration page from the HTML template | Implemented |
| RegisterHandler | [handlers/auth.go](../handlers/auth.go) | Receives POST data, hashes the password, inserts user into SQLite | Implemented |
| RegisterDispatcher | [handlers/auth.go](../handlers/auth.go) | Chooses between GET and POST handling for /register | Implemented |
| HomeHandler | [handlers/home.go](../handlers/home.go) | Renders the home page from the HTML template | Implemented |
| InitDB | [database/db.go](../database/db.go) | Opens SQLite, creates tables, stores the global DB handle | Implemented |
| InitDB | [internal/database/db.go](../internal/database/db.go) | Alternative database initialization implementation | Implemented but not used |

---

## 6. Current implementation details

### 6.1 Main entry point

File: [main.go](../main.go)

What it does:

- imports the database package and handlers package
- calls database.InitDB() before the server starts
- registers routes for the homepage and registration flow
- starts the server on localhost:8080

### 6.2 Home page

File: [handlers/home.go](../handlers/home.go)

What it does:

- loads [ui/index.html](../ui/index.html)
- serves it as the homepage

Current status:

- very basic template
- no data is loaded from the database yet

### 6.3 Registration flow

File: [handlers/auth.go](../handlers/auth.go)

What it does:

- shows the registration form on GET /register
- accepts username, email, and password on POST /register
- hashes the password using bcrypt
- inserts a new row into the users table
- returns a success or conflict response

Current status:

- registration page works
- user registration into SQLite works
- there is no login flow yet
- there is no validation beyond basic form submission

### 6.4 Database layer

File: [database/db.go](../database/db.go)

What it does:

- opens the SQLite file at ./forum.db
- creates the following tables if they do not exist:
  - users
  - categories
  - posts
  - comments
  - likes
  - sessions

Current status:

- tables are defined in SQL
- the app currently only uses the users table in practice
- the rest of the schema is planned but not yet implemented in handlers

### 6.5 Templates

File: [ui/register.html](../ui/register.html)

- contains a simple form with username, email, and password fields
- posts to /register

File: [ui/index.html](../ui/index.html)

- contains a basic welcome page

---

## 7. What is already implemented

The following pieces are already in place:

- Go HTTP server starts correctly
- SQLite database connection is initialized
- users table exists
- registration page is available
- registration handler accepts form input
- passwords are hashed before storage
- basic templates are being served

---

## 8. What is still missing

The project is not finished yet. The following features are still missing or incomplete:

### Authentication

- login handler
- logout handler
- session creation
- session validation
- cookie-based authentication
- protected routes for authenticated users

### Forum content

- creating posts
- displaying posts
- creating comments
- displaying comments

### Interaction features

- likes and dislikes for posts/comments
- visible counts for likes/dislikes

### Organization features

- categories
- filtering by category
- filtering by created posts
- filtering by liked posts

### Project quality

- Docker setup
- tests
- better error handling and templates
- code organization cleanup

---

## 9. Where the project is stopping right now

The current stopping point is the start of the user account flow.

In plain terms:

- the app can run
- the database can be initialized
- the registration page exists
- a user can be inserted into the database

What is not implemented yet:

- login
- sessions
- posting content
- comments
- likes
- user-specific pages and filters

So the next logical milestone is:

1. implement login and session cookies
2. protect forum actions behind authentication
3. add post creation and display

---

## 10. Recommended next implementation order

To continue this project in a clean way, this is the best order:

1. Standardize the database package
   - decide whether the app will use [database/db.go](../database/db.go) or [internal/database/db.go](../internal/database/db.go)
   - avoid having two similar implementations

2. Implement login/logout
   - create a login page
   - validate credentials against the users table
   - create a session row in the sessions table
   - set a cookie

3. Implement post creation and display
   - add a new page for creating posts
   - save posts to the posts table
   - show posts on the homepage or a dedicated posts page

4. Implement comments
   - add comments linked to posts
   - display them under each post

5. Implement likes/dislikes
   - add buttons for like/dislike
   - update the likes table

6. Implement categories and filters
   - allow posts to belong to one or more categories
   - add filtering by category, own posts, and liked posts

7. Add Docker support
   - create a Dockerfile if missing
   - document how to run the app in a container

8. Add tests
   - cover registration logic
   - cover database initialization
   - cover core handlers

---

## 11. Important notes for the next developer or AI assistant

### Duplicate database packages

There are two database packages in the project:

- [database/db.go](../database/db.go)
- [internal/database/db.go](../internal/database/db.go)

The app currently uses [database/db.go](../database/db.go) because [main.go](../main.go) imports it. The internal version is currently redundant and should be cleaned up later.

### Templates and working directory

The templates are loaded using relative paths such as ui/index.html and ui/register.html. This means the app should be run from the project root so those paths resolve correctly.

### Database file location

The database is opened using ./forum.db. That means the SQLite file is created in the current working directory.

### Windows setup note

If you are running this on Windows, you may need the setup described in [docs/fix.txt](fix.txt) so that SQLite works correctly with Go and CGO.

---

## 12. Quick start commands

From the project root:

```powershell
go run .
```

Then open:

```text
http://localhost:8080/
```

Registration is available at:

```text
http://localhost:8080/register
```

---

## 13. Handoff summary

This repository is currently at the registration milestone.

If you continue the project, the next meaningful feature to implement is authentication with sessions. Once that works, the app can move into posts, comments, likes, and filters.

This document should be used as the main reference point whenever the project is resumed.
