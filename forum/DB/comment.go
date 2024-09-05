package DB

import (
	"database/sql"
	"log"
	"fmt"
    "net/http"
)


type Comment struct {
    CommentID int    `json:"comment_id"`
    PostID    int    `json:"post_id"`
    Author    string `json:"author"`
    Content   string `json:"content"`
    Date      string `json:"date"`
    Likes     int    `json:"likes"`
    Dislikes  int    `json:"dislikes"`
    
}



func CreateCommentTable() {
	// Open the database connection
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	//Sql to create table
	createTableCommentSQL := `
	CREATE TABLE IF NOT EXISTS comments (
    comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    post_id INTEGER,
    content TEXT NOT NULL,
    comment_date TEXT NOT NULL,
    num_likes INTEGER DEFAULT 0,
    num_dislikes INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);`

	_, err = db.Exec(createTableCommentSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func GetCommentsByPostID(postID int) (*[]Comment, error) {
    query := `
    SELECT 
        c.comment_id, 
        c.post_id, 
        u.username, 
        c.content, 
        c.comment_date, 
        c.num_likes, 
        c.num_dislikes
    FROM comments c
    JOIN users u ON c.user_id = u.id
    WHERE c.post_id = ?`

    rows, err := DB.Query(query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []Comment
    for rows.Next() {
        var comment Comment
        err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Author, &comment.Content, &comment.Date, &comment.Likes, &comment.Dislikes)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return &comments, nil
}


func AddComment(comment Comment) (bool, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return false, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `INSERT INTO comments (post_id, user_id, content, comment_date, num_likes, num_dislikes) 
              VALUES (?, ?, ?, ?, ?, ?)`

    stmt, err := db.Prepare(query)
    if err != nil {
        return false, err
    }
    defer stmt.Close()

    _, err = stmt.Exec(comment.PostID, comment.Author, comment.Content, comment.Date, comment.Likes, comment.Dislikes)
    if err != nil {
        return false, err
    }

    return true, nil
}

func GetUserCommentAction(userID, commentID int) (string, error) {
    var actionType string
    query := `SELECT action_type FROM user_comments WHERE user_id = ? AND comment_id = ?`
    err := DB.QueryRow(query, userID, commentID).Scan(&actionType)
    if err == sql.ErrNoRows {
        return "", nil // No prior action
    }
    return actionType, err
}

func IncrementLikeC(userID, commentID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }


    // Increment the like count
    _, err = tx.Exec(`UPDATE comments SET num_likes = num_likes + 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE comments SET num_dislikes = num_dislikes - 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_comments SET action_type = 'like'  WHERE user_id = ? AND comment_id = ?`,userID, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func IncrementLikeC12(userID, commentID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE comments SET num_likes = num_likes - 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_comments SET action_type = ''  WHERE user_id = ? AND comment_id = ?`, userID,commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}


func IncrementDisLikeC12(userID, commentID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE comments SET num_dislikes = num_dislikes - 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_comments SET action_type = ''  WHERE user_id = ? AND comment_id = ?`,userID,  commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}


func IncrementLikeC1(userID, commentID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE comments SET num_likes = num_likes + 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_comments SET action_type = 'like'  WHERE user_id = ? AND comment_id = ?`, userID,commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func IncrementDisLikeC1(userID, commentID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE comments SET num_dislikes = num_dislikes + 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_comments SET action_type = 'dislike'  WHERE user_id = ? AND comment_id = ?`,userID,  commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}


func IncrementDisLikeC(userID, commentID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE comments SET num_dislikes = num_dislikes + 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE comments SET num_likes = num_likes - 1 WHERE comment_id = ?`, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_comments SET action_type = 'dislike'  WHERE user_id = ? AND comment_id = ?`,userID, commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}


func InsertCommentDb(w http.ResponseWriter ,postID int, userID int, content string, formattedDate string){
    _, err := DB.Exec("INSERT INTO comments (post_id, user_id, content, comment_date) VALUES (?, ?, ?, ?)", postID, userID, content, formattedDate)
    if err != nil {
        http.Error(w, "Failed to submit comment", http.StatusInternalServerError)
        return
    }

}

func InsertFirstAction(w http.ResponseWriter, userID int, commentID int, action string, existingAction string){
    
    if existingAction == "" {
        _, err := DB.Exec(`INSERT INTO user_comments (user_id, comment_id, action_type) VALUES (?, ?, ?)`, userID, commentID, action)
            if err != nil {
                return
            }

    }

}

