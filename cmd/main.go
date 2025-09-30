package main

import (
	"auth-api/config"
	"auth-api/db"
	api "auth-api/handlers"
	mw "auth-api/middleware"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/health", mw.Logger(health_handler))

	http.HandleFunc("/login", mw.Logger(api.LoginHandler))
	http.HandleFunc("/register", mw.Logger(api.RegisterHandler))

	http.HandleFunc("/secret", mw.Logger(mw.CheckJwt(api.SecretHandler)))

	//err := db.InitDB("token_master", "jwt_users", "tokenPass", "auth-db")		// host name comes from docker-compose.yml
	err := db.InitDB(config.User, config.Dbname, config.Password, config.Host)
	if err != nil {
		log.Printf("failed initilizing the db: %v", err)
	}

	http.ListenAndServe(":8976", nil)

	/*
	   nil is the multiplexer...which is kinda like a switchboard..
	   a multiplexer, sees the path, and call the specific handler function
	   calling handleFunc add a rule to the multiplexer.
	   there is a default multiplexer, and you can write a custom multiplexer
	*/
}

func health_handler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("<h1>ALIVE AND WELL...ish</h1>\n"))
}
