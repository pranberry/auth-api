package main

import (
	"jwt-auth/config"
	"jwt-auth/db"
	"jwt-auth/secret"
	"jwt-auth/user"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func main() {

	mux := mux.NewRouter()
	
	mux.HandleFunc("/health", health_handler).Methods("GET")
	mux.HandleFunc("/login", user.LoginHandler).Methods("POST")
	mux.HandleFunc("/register", user.RegisterHandler).Methods("POST")
	mux.HandleFunc("/secret", secret.SecretHandler).Methods("GET")

	// Initialize the DB
	err := db.InitDB(config.User, config.Dbname, config.Password, config.Host)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	// Starting listening on 8080
	wg.Add(1)
	go http.ListenAndServe(":9000", mux)
	
	wg.Wait()
}

func health_handler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("<h1>ALIVE AND WELL...ish</h1>"))
}
