package DB

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const DBPath string = "./DB/forum.db"

var DB *sql.DB

func CreateDB() {
	//Open the DB file. Will create a DB if it does not exist
	var err error
	DB, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	//Create the DB tables 
	CreateUsersTable()
	CreateSessionsTable()
	CreatePostsTable()
	CreateCommentTable()
	CreateUserInterTable()
	CreateUserInterCommentTable() 
}
