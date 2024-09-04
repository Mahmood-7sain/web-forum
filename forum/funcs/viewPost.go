package funcs

import (
	"fmt"
	"forum/DB"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostDetails struct {
	Post     DB.Post
	Comments []DB.Comment
}

func ViewPost(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query string
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	// Fetch the post details by ID
	post, err := DB.GetPostByID(postID)
	if post == nil || err != nil{
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	// Fetch the comments for the post
	comments, err := DB.GetCommentsByPostID(postID)
	if err != nil {
		res := Error{}
		res.Code = 500
		res.Status = "Internal Server Error"
		ErrorHandler(w, r, &res)
		return
	}

	// Create the PostDetails struct
	postDetails := PostDetails{
		Post:     *post,
		Comments: *comments,
	}
	//     log.Println("Post ID:", postID)
	// log.Println("Fetched Comments:", comments)

	// Render the template with the post and comment data
	err = tpl.ExecuteTemplate(w, "single-post.html", postDetails)
	if err != nil {
		res := Error{}
		res.Code = 500
		res.Status = "Internal Server Error"
		ErrorHandler(w, r, &res)
		return
		fmt.Println("Template execution error:", err)
	}
}

func HandleCommentSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	// Get the logged-in user's ID from the session or context
	userID, loggedIn := GetSession(r)
	if !loggedIn {
		res := Error{}
		res.Code = 403
		res.Status = "Unauthorized"
		ErrorHandler(w, r, &res)
		return
	}

	// Get the post ID and comment content from the form
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")

	// Convert post ID to an integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	if strings.ReplaceAll(content, " ", "") == "" {
		// Redirect back to the single post page
		http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
		return
	}

	// Get the current date and time
	formattedDate := time.Now().Format("2006-01-02 15:04:05")

	DB.InsertCommentDb(w, postID, userID, content, formattedDate)

	// Redirect back to the single post page
	http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
}

func HandelCommmentLike(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	// Get the logged-in user's ID from the session or context
	userID, loggedIn := GetSession(r)
	if !loggedIn {
		res := Error{}
		res.Code = 403
		res.Status = "Unauthorized"
		ErrorHandler(w, r, &res)
		return
	}

	// Get the post ID and comment content from the form
	postIDStr := r.FormValue("post_id")
	//content := r.FormValue("content")

	// Convert post ID to an integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	action := r.FormValue("action")
	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}
	existingAction, err := DB.GetUserCommentAction(userID, commentID)
	if err != nil {
		res := Error{}
		res.Code = 500
		res.Status = "Internal Server Error"
		ErrorHandler(w, r, &res)
		return
	}

	DB.InsertFirstAction(w, userID, commentID, action, existingAction)

	if action == "like" && existingAction == "like" {
		// If the user already liked the post, do nothing or display a message
		http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
		return
	}
	if action == "like" && existingAction == "dislike" {
		// If the user had previously disliked, decrease the dislike and increase the like
		err = DB.IncrementLikeC(userID, commentID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}
	}
	if action == "dislike" && existingAction == "dislike" {
		http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
		return
		// If no prior action, just increment the like

	}
	if action == "dislike" && existingAction == "like" {
		err = DB.IncrementDisLikeC(userID, commentID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

	}
	if action == "like" && existingAction == "" {
		err = DB.IncrementLikeC1(userID, commentID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

	}
	if action == "dislike" && existingAction == "" {

		err = DB.IncrementDisLikeC1(userID, commentID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}
	}

	// Redirect back to the single post page
	http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)

}

func HandleLike(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}

	userID, loggedIn := GetSession(r)
	if !loggedIn {
		res := Error{}
		res.Code = 403
		res.Status = "Unauthorized"
		ErrorHandler(w, r, &res)
		return
	}

	// Assuming you have a way to get the logged-in user's ID
	// here is the actual action that will happen
	action := r.FormValue("action")
	// Check if the user already liked or disliked this post
	existingAction, err := DB.GetUserAction(userID, postID)
	if err != nil {
		res := Error{}
		res.Code = 500
		res.Status = "Internal Server Error"
		ErrorHandler(w, r, &res)
		return
	}

	DB.InsertIIntoUserInter(w, userID, postID, action, existingAction)

	if action == "like" && existingAction == "like" {
		// If the user already liked the post, do nothing or display a message
		http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
		return
	}
	if action == "like" && existingAction == "dislike" {
		// If the user had previously disliked, decrease the dislike and increase the like
		err = DB.IncrementLike(userID, postID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}
	}
	if action == "dislike" && existingAction == "dislike" {
		http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
		return
		// If no prior action, just increment the like

	}
	if action == "dislike" && existingAction == "like" {
		err = DB.IncrementDisLike(userID, postID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

	}
	if action == "like" && existingAction == "" {
		err = DB.IncrementLike1(userID, postID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

	}
	if action == "dislike" && existingAction == "" {

		err = DB.IncrementDisLike1(userID, postID)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/viewSinglePost?id=%d", postID), http.StatusSeeOther)
}
