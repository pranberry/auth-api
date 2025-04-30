package main

import (
	"net/http"
	"jwt-auth/user"
	"jwt-auth/secret"
)

func main() {
	// HandleFunc is from the http lib
	// associates a URL path with a function
	http.HandleFunc("/health", health_handler)

	http.HandleFunc("/login", user.LoginHandler)
	http.HandleFunc("/register", user.RegisterHandler)
	
	http.HandleFunc("/secret", secret.SecretHandler)

	//listen on port 8080...blocking call
	http.ListenAndServe(":8080", nil)
	/*
	   nil is the multiplexer...which is kinda like a switchboard.
	   a multiplexer, see the path, and call the specifiec handler function
	   calling handleFunc add a rule to the multiplexer.
	   there is a default multiplexer, and you can write a custom multiplexer
	*/ 
}


func health_handler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("<h1>ALIVE AND WELL...ish</h1>"))
}


/*
    TO-DO
    - change all error messages to send back JSON, not plaintext
    - eventually, add https to the mix, so we're getting encrypted data
*/

