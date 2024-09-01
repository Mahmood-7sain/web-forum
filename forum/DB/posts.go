package DB

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	//"encoding/json"
	
)

type Post struct {
    ID         int      `json:"id"`
    UserID     int      `json:"user_id"`
    UserName   string   `json:"user_name"`
    Title      string   `json:"title"`
    Content    string   `json:"content"`
    Date       string   `json:"date"`
    Categories []string `json:"categories"`
    Likes      int      `json:"likes"`
    Dislikes   int      `json:"dislikes"`
    Comments   int `json:"comments"`
   NComments []Comment   `json:"comment"`
}

type PostDetails struct {
    post1 Post
    comm []Comment
}


func CreatePostsTable() {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the posts table
	createPostsTableSQL := `CREATE TABLE IF NOT EXISTS posts (
        post_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
        post_title TEXT NOT NULL,
        post_content TEXT NOT NULL,
        post_date TEXT NOT NULL,
        category TEXT NOT NULL,
        num_likes INTEGER DEFAULT 0,
        num_dislikes INTEGER DEFAULT 0,
        num_comments INTEGER DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
	_, err = db.Exec(createPostsTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func AddPost(post Post) (bool, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return false, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    // Convert the Categories slice to a comma-separated string
    categories := strings.Join(post.Categories, ",")

    query := `INSERT INTO posts (user_id, username, post_title, post_content, post_date, category, num_likes, num_dislikes, num_comments) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

    stmt, err := db.Prepare(query)
    if err != nil {
        return false, err
    }
    defer stmt.Close()

    _, err = stmt.Exec(post.UserID, post.UserName, post.Title, post.Content, post.Date, categories, post.Likes, post.Dislikes, post.Comments)
    if err != nil {
        return false, err
    }

    return true, nil
}

func FetchPosts() ([]Post, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT post_id, user_id, username, post_title, post_content, post_date, category, num_likes, num_dislikes, num_comments FROM posts")
    if err != nil {
        return nil, fmt.Errorf("failed to fetch posts: %v", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var categories string

        err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Title, &post.Content, &post.Date, &categories, &post.Likes, &post.Dislikes, &post.Comments)
        if err != nil {
            return nil, fmt.Errorf("failed to scan post data: %v", err)
        }

        // Split the categories by comma into a slice of strings
        post.Categories = strings.Split(categories, ",")

        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %v", err)
    }

    return posts, nil
}

func GetPostByID(postID int) (*Post, error) {
    db, err := InitDB()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    var post Post
    var categoryCSV string

    // Join the posts and users tables to fetch the username along with the post details
    query := `
        SELECT p.post_id, p.user_id, u.username, p.post_title, p.post_content, p.post_date, p.category, p.num_likes, p.num_dislikes, p.num_comments 
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.post_id = ?`

    err = db.QueryRow(query, postID).Scan(
        &post.ID,
        &post.UserID,
        &post.UserName,
        &post.Title,
        &post.Content,
        &post.Date,
        &categoryCSV,
        &post.Likes,
        &post.Dislikes,
        &post.Comments,
    )
    if err != nil {
        return nil, err
    }

    // Convert the comma-separated string into a slice of strings
    post.Categories = strings.Split(categoryCSV, ",")

    return &post, nil
}


func FetchUserPosts(id int) ([]Post, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := "SELECT post_id, user_id, username, post_title, post_content, post_date, category, num_likes, num_dislikes, num_comments FROM posts WHERE user_id=?"

    rows, err := db.Query(query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch posts: %v", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var categories string

        err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Title, &post.Content, &post.Date, &categories, &post.Likes, &post.Dislikes, &post.Comments)
        if err != nil {
            return nil, fmt.Errorf("failed to scan post data: %v", err)
        }

        // Split the categories by comma into a slice of strings
        post.Categories = strings.Split(categories, ",")

        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %v", err)
    }

    return posts, nil
}


func FetchLikedPostsByUserID(userID int) ([]Post, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `
        SELECT posts.post_id, posts.user_id, posts.username, posts.post_title, posts.post_content, posts.post_date, posts.category, posts.num_likes, posts.num_dislikes, posts.num_comments
        FROM posts
        INNER JOIN user_interactions ON posts.post_id = user_interactions.post_id
        WHERE user_interactions.user_id = ? AND user_interactions.action_type = 'like'
    `

    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch liked posts: %v", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var categories string

        err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Title, &post.Content, &post.Date, &categories, &post.Likes, &post.Dislikes, &post.Comments)
        if err != nil {
            return nil, fmt.Errorf("failed to scan liked post data: %v", err)
        }

        // Split the categories by comma into a slice of strings
        post.Categories = strings.Split(categories, ",")

        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating liked posts rows: %v", err)
    }

    return posts, nil
}




func GetPostsByCategory(cat string) ([]Post, error) {
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `
        SELECT post_id, user_id, username, post_title, post_content, post_date, category, num_likes, num_dislikes, num_comments
        FROM posts
        WHERE category LIKE ?
    `

    rows, err := db.Query(query, "%"+cat+"%")
    if err != nil {
        return nil, fmt.Errorf("failed to fetch posts by category: %v", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var categories string

        err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Title, &post.Content, &post.Date, &categories, &post.Likes, &post.Dislikes, &post.Comments)
        if err != nil {
            return nil, fmt.Errorf("failed to scan post data: %v", err)
        }

        // Split the categories by comma into a slice of strings
        post.Categories = strings.Split(categories, ",")

        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating posts rows: %v", err)
    }

    return posts, nil
}
