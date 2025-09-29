package main

import (
	"jwt-auth/config"
	"jwt-auth/db"
	"jwt-auth/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/health", health_handler)

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)

	http.HandleFunc("/secret", handlers.SecretHandler)

	//err := db.InitDB("token_master", "jwt_users", "tokenPass", "auth-db")		// host name comes from docker-compose.yml
	err := db.InitDB(config.User, config.Dbname, config.Password, config.Host)
	if err != nil {
		log.Fatal("died initilizing the db: ", err)
	}

	go http.ListenAndServe(":8976", nil)
	select {}

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
