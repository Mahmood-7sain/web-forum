package main

import (
	"fmt"
	"forum/DB"
	"forum/funcs"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func main() {
	var err error
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	funcs.SetT(tpl)
	DB.CreateDB()

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	http.HandleFunc("/", funcs.HomeHandler)
	http.HandleFunc("/login", funcs.Login)
	http.HandleFunc("/logout", funcs.LogoutHandler)
	http.HandleFunc("/signup", funcs.Signup)
	http.HandleFunc("/createPost", funcs.SessionMiddleware(funcs.CreateHandler))
	http.HandleFunc("/viewSinglePost", funcs.ViewPost)
	http.HandleFunc("/comment", funcs.SessionMiddleware(funcs.HandleCommentSubmission))
	http.HandleFunc("/like-dislike", funcs.SessionMiddleware(funcs.HandleLike))
	http.HandleFunc("/comment-like", funcs.SessionMiddleware(funcs.HandelCommmentLike))

	fmt.Println("Starting Server...")
	fmt.Println("Listening on port 8080")
	fmt.Println("Navigate to: http://localhost:8080")
	ServerError := http.ListenAndServe(":8080", nil)
	if ServerError != nil {
		log.Fatal(ServerError)
	}
}
