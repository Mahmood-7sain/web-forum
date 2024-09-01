package DB

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// DBPath is the path to the SQLite database file.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}
// CreateSessionsTable creates the sessions table in the database.
func CreateSessionsTable() {
	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createSessionsTableSQL := `CREATE TABLE IF NOT EXISTS sessions (
        session_id TEXT PRIMARY KEY,
        expiration DATETIME,
        user_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	_, err = db.Exec(createSessionsTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
