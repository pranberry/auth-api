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

	http.HandleFunc("/health", mw.Logger(api.HealthHandler))

	http.HandleFunc("/login", mw.Logger(api.LoginHandler))
	http.HandleFunc("/register", mw.Logger(api.RegisterHandler))

	http.HandleFunc("/secret", mw.Logger(mw.CheckJwt(api.SecretHandler)))

	// Initialize the DB. All these values live in the .env or .env.local
	err := db.InitDB(config.User, config.DbName, config.Password, config.Host)

	if err != nil {
		log.Fatalf("failed initializing the db: %v", err)
	}

	http.ListenAndServe(":8976", nil)

}
