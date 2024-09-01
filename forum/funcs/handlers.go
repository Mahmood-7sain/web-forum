package funcs

import (
	"fmt"
	"forum/DB"
	"log"
	"net/http"
	"os"
)

type Result struct {
	Posts  []DB.Post
	Logged bool
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			res := Error{}
			res.Code = 404
			res.Status = "Not Found"
			ErrorHandler(w, r, &res)
			return
		}

		result := Result{}

		// Get the post ID from the query string
		category := r.URL.Query().Get("type")

		id, isLogged := GetSession(r)

		if category == "" || category == "all" {
			posts, errFetch := DB.FetchPosts() // FetchPosts should return a slice of all Post objects
			if errFetch != nil {
				log.Printf("Error fetching posts: %v", errFetch)
				res := Error{}
				res.Code = 500
				res.Status = "Internal Server Error"
				ErrorHandler(w, r, &res)
				return
			}

			result.Posts = posts
			result.Logged = isLogged

			err := tpl.ExecuteTemplate(w, "index.html", &result)
			if err != nil {
				log.Printf("Error rendering template: %v", err)
				res := Error{}
				res.Code = 500
				res.Status = "Internal Server Error"
				ErrorHandler(w, r, &res)
				return
			}
		} else if category == "my-posts" {
			if id == 0 || !isLogged {
				// Redirect the user to the login page
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				posts, err := DB.FetchUserPosts(id)
				if err != nil {
					res := Error{}
					res.Code = 500
					res.Status = "Internal Server Error"
					ErrorHandler(w, r, &res)
					return
				}

				result := Result{}
				result.Posts = posts
				result.Logged = isLogged
				err = tpl.ExecuteTemplate(w, "index.html", &result)
				if err != nil {
					log.Printf("Error rendering template: %v", err)
					res := Error{}
					res.Code = 500
					res.Status = "Internal Server Error"
					ErrorHandler(w, r, &res)
					return
				}

			}
		} else if category == "liked" {
			if id == 0 || !isLogged {
				// Redirect the user to the home page
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				posts, err := DB.FetchLikedPostsByUserID(id)
				if err != nil {
					fmt.Println(err)
					res := Error{}
					res.Code = 500
					res.Status = "Internal Server Error"
					ErrorHandler(w, r, &res)
					return
				}

				result := Result{}
				result.Posts = posts
				result.Logged = isLogged
				err = tpl.ExecuteTemplate(w, "index.html", &result)
				if err != nil {
					log.Printf("Error rendering template: %v", err)
					res := Error{}
					res.Code = 500
					res.Status = "Internal Server Error"
					ErrorHandler(w, r, &res)
					return
				}
			}
		} else if category == "tech" || category == "books" || category == "sports" {
			//Get the posts by category
			posts, err := DB.GetPostsByCategory(category)
			if err != nil {
				fmt.Println(err)
				res := Error{}
				res.Code = 500
				res.Status = "Internal Server Error"
				ErrorHandler(w, r, &res)
				return
			}

			result := Result{}
			result.Posts = posts
			result.Logged = isLogged
			err = tpl.ExecuteTemplate(w, "index.html", &result)
			if err != nil {
				log.Printf("Error rendering template: %v", err)
				res := Error{}
				res.Code = 500
				res.Status = "Internal Server Error"
				ErrorHandler(w, r, &res)
				return
			}
		} else {
			res := Error{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w, r, &res)
			return
		}
	} else {
		res := Error{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearSession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Function to handle any errors. Receives a code and status and executes the error.html file with the corresponding status and status code
func ErrorHandler(w http.ResponseWriter, r *http.Request, res *Error) {
	w.WriteHeader(res.Code)
	err := tpl.ExecuteTemplate(w, "error.html", res)
	if err != nil {
		fmt.Println("Error with error.html")
		os.Exit(2)
	}
}
