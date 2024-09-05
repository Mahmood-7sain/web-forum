package funcs

import (
	"fmt"
	"forum/DB"
	"net/http"
	"strings"
	"time"
)

type CreateData struct {
	Username string
	Date     string
}

type ErrorMsg struct {
	EmptyField string
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, found := GetSession(r)
		if id == 0 && !found {
			fmt.Println("No session found or id 0")
			res := Error{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w, r, &res)
			return
		} else {
			username, errGetUsername := DB.GetUsername(id)
			if errGetUsername != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			} else if username == "" {
				fmt.Println("No username")
				res := Error{}
				res.Code = 400
				res.Status = "Bad Request"
				ErrorHandler(w, r, &res)
				return
			}

			data := CreateData{}
			data.Username = username

			// Get the current date and time
			now := time.Now()

			// Format the date as dd/mm/yy
			formattedDate := now.Format("02/01/06")

			data.Date = formattedDate

			err := tpl.ExecuteTemplate(w, "create-post.html", nil)
			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

		}
	} else if r.Method == http.MethodPost {
		//If the method is post then get the form data and validate it
		title := r.FormValue("post-title")
		categories := r.Form["category"]
		content := r.FormValue("post-content")

		if strings.ReplaceAll(title, " ", "") == "" || len(categories) == 0 || strings.ReplaceAll(content, " ", "") == "" {
			// Redirect the user back to same page if he send empty form
			msg := ErrorMsg{}
			msg.EmptyField = "Please fill all the fields!"

			err := tpl.ExecuteTemplate(w, "create-post.html", &msg)

			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

		} else {

			id, found := GetSession(r)
			if id == 0 && !found {
				fmt.Println("No session found or id 0")
				res := Error{}
				res.Code = 400
				res.Status = "Bad Request"
				ErrorHandler(w, r, &res)
				return
			}

			userExists, err := DB.CheckUser(id)
			if err != nil {
				fmt.Println("2")
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

			if !userExists {
				res := Error{}
				res.Code = 400
				res.Status = "Bad Request"
				ErrorHandler(w, r, &res)
				return
			}

			username, errGetUsername := DB.GetUsername(id)
			if errGetUsername != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

			if username == "" {
				res := Error{}
				res.Code = 400
				res.Status = "Bad Request"
				ErrorHandler(w, r, &res)
				return
			}

			// Get the current date and time
			now := time.Now()

			// Format the date as dd/mm/yy
			formattedDate := now.Format("02/01/06")

			newPost := DB.Post{}
			newPost.UserID = id
			newPost.UserName = username
			newPost.Title = title
			newPost.Content = content
			newPost.Date = formattedDate
			newPost.Categories = categories
			newPost.Likes = 0
			newPost.Dislikes = 0
			newPost.Comments = 0

			success, errAdd := DB.AddPost(newPost)
			if errAdd != nil || !success {
				fmt.Println(errAdd)
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

			// Redirect the user to the home page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
}
