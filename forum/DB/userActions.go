package DB

import (
	"database/sql"
	"log"
	"net/http"
)

func CreateUserInterTable() {
	// Open the database connection
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	//Sql to create table
	createTableCommentSQL := `
	CREATE TABLE IF NOT EXISTS user_interactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    post_id INTEGER,
    action_type TEXT, -- 'like' or 'dislike'
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);`

	_, err = db.Exec(createTableCommentSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUserInterCommentTable() {
	// Open the database connection
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	//Sql to create table
	createTableCommentSQL := `
	CREATE TABLE IF NOT EXISTS user_comments (
    user_id INTEGER,
    comment_id INTEGER,
    action_type TEXT, -- 'like' or 'dislike'
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id)
);`

	_, err = db.Exec(createTableCommentSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertIIntoUserInter(w http.ResponseWriter, userID int, postID int, action string, existingAction string) {
	if existingAction == "" {
		_, err := DB.Exec(`INSERT INTO user_interactions (user_id, post_id, action_type) VALUES (?, ?, ?)`, userID, postID, action)
		if err != nil {
			http.Error(w, "Unable to record interaction", http.StatusInternalServerError)
			return
		}

	}

}
