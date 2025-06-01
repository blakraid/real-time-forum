package main

import (
	"fmt"
	"log"
	"rtf/database"
	"rtf/server"
	"net/http"
)



func main(){
	port := "8167"
	if err := database.CreateDatabase(); err != nil {
		log.Fatal("Problem in Create database")
	}

	// File server
	FileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/",http.StripPrefix("/static",FileServer))

	// Content Home Page

	http.HandleFunc("/show_posts", server.ShowPosts)
	http.HandleFunc("/post_submit", server.PostSubmit)
	http.HandleFunc("/get_categories", server.GetCategories)
	http.HandleFunc("/comment_submit", server.CommentSubmit)
	http.HandleFunc("/interact", server.HandleInteract)
	
	// authentication handler
	http.HandleFunc("/regester",server.RegesterHandler)
	http.HandleFunc("/",server.HomeHandler)
	http.HandleFunc("/login",server.LoginHandler)
	http.HandleFunc("/check-session", server.CheckSessionHandler)
	http.HandleFunc("/logout", server.LogoutHandler)

	http.HandleFunc("/ws", server.HandleConnection)
	http.HandleFunc("/users", server.HandleUsers)
	http.HandleFunc("/messages", server.GetMessages)


	// Running Server
	fmt.Println("http://localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port,nil))

}