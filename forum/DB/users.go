package DB

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	UserID       int
	UserName     string
	UserEmail    string
	UserPassword string
}

func CreateUsersTable() {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table
	createTableSQLUsers := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQLUsers)
	if err != nil {
		log.Fatal(err)
	}
}

// CheckEmailExists checks if the given email already exists in the users table.
func CheckEmailExists(email string) (bool, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return false, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=? LIMIT 1)`
	err = db.QueryRow(query, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}

	return exists, nil
}

// CheckEmailExists checks if the given email already exists in the users table.
func CheckUserExists(username string) (bool, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return false, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username=? LIMIT 1)`
	err = db.QueryRow(query, username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}

	return exists, nil
}

func AddUser(user User) error {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the SQL statement with the user data
	_, err = stmt.Exec(user.UserName, user.UserEmail, user.UserPassword)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// Returns the user with the passed email, error if not found
func GetUser(email string) (*User, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	user := &User{}

	query := "SELECT id, username, email, password FROM users WHERE email = ?"

	row := db.QueryRow(query, email)
	errScan := row.Scan(&user.UserID, &user.UserName, &user.UserEmail, &user.UserPassword)
	if errScan != nil {
		if errScan == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with email %s", email)
		}
		return nil, err
	}

	return user, nil
}

func GetUsername(id int) (string, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return "", fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	var username string

	query := "SELECT username FROM users WHERE id = ?"

	row := db.QueryRow(query, id)
	errScan := row.Scan(&username)
	if errScan != nil {
		if errScan == sql.ErrNoRows {
			return "", fmt.Errorf("no user with id %s", id)
		}
		return "", err
	}

	return username, nil
}

func CheckUser(id int) (bool, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return false, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	var exists bool

	query := "SELECT COUNT(*) FROM users WHERE id = ?"
	errScan := db.QueryRow(query, id).Scan(&exists)
	if errScan != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}
