package main

import (
	"jwt-auth/db"
	"jwt-auth/secret"
	"jwt-auth/user"
	"log"
	"net/http"
)

func main() {
	// HandleFunc is from the http lib
	// associates a URL path with a function
	http.HandleFunc("/health", health_handler)

	http.HandleFunc("/login", user.LoginHandler)
	http.HandleFunc("/register", user.RegisterHandler)
	
	http.HandleFunc("/secret", secret.SecretHandler)

	err := db.InitDB("token_master", "jwt_users")
	if err != nil{
		log.Fatal("died initilizing the db: ", err)
	}

	//listen on port 8080...blocking call
	go http.ListenAndServe(":8080", nil)
	go http.ListenAndServe(":8081", nil)
	select{}

	/*
	   nil is the multiplexer...which is kinda like a switchboard.
	   a multiplexer, sees the path, and call the specifiec handler function
	   calling handleFunc add a rule to the multiplexer.
	   there is a default multiplexer, and you can write a custom multiplexer
	*/ 
}

func health_handler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("<h1>ALIVE AND WELL...ish</h1>"))
}