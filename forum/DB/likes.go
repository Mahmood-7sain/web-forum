package DB

import (
	"database/sql"
	//"log"
)


func DecrementDisLike1(userID, postID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE posts SET num_dislikes = num_dislikes - 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_interactions SET action_type = ''  WHERE user_id = ? AND post_id = ?`,userID,  postID)
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


func DecrementLike1(userID, postID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE posts SET num_likes = num_likes - 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_interactions SET action_type = ''  WHERE user_id = ? AND post_id = ?`, userID,postID)
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

func IncrementLike(userID, postID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }


    // Increment the like count
    _, err = tx.Exec(`UPDATE posts SET num_likes = num_likes + 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE posts SET num_dislikes = num_dislikes - 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_interactions SET action_type = 'like'  WHERE user_id = ? AND post_id = ?`,userID, postID)
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


func IncrementLike1(userID, postID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE posts SET num_likes = num_likes + 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_interactions SET action_type = 'like'  WHERE user_id = ? AND post_id = ?`, userID,postID)
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

func IncrementDisLike1(userID, postID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE posts SET num_dislikes = num_dislikes + 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_interactions SET action_type = 'dislike'  WHERE user_id = ? AND post_id = ?`,userID,  postID)
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


func IncrementDisLike(userID, postID int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    // Increment the like count
    _, err = tx.Exec(`UPDATE posts SET num_dislikes = num_dislikes + 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE posts SET num_likes = num_likes - 1 WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE user_interactions SET action_type = 'dislike'  WHERE user_id = ? AND post_id = ?`,userID, postID)
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



func GetUserAction(userID, postID int) (string, error) {
    var actionType string
    query := `SELECT action_type FROM user_interactions WHERE user_id = ? AND post_id = ?`
    err := DB.QueryRow(query, userID, postID).Scan(&actionType)
    if err == sql.ErrNoRows {
        return "", nil // No prior action
    }
    return actionType, err
}

