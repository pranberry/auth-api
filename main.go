package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type service_user struct {
    // in GO, field names should start with capital letters to be unmarshaled (decoded from JSON)
    User_Name string    `json:"username"`
    // the bits in the back-tics are "struct-tags", this tells json.decode() what to look for
    Password string     `json:"password"`
    Location string
    IP_addr string
}

func main() {
	// HandleFunc is from the http lib
	// associates a URL path with a function
	http.HandleFunc("/health", health_handler)
	http.HandleFunc("/login", login_handler)
	http.HandleFunc("/register", register_handler)

	//listen on port 8080...blocking call
	http.ListenAndServe(":8080", nil)
	/*
	   nil is the multiplexer...which is kinda like a switchboard.
	   a multiplexer, see the path, and call the specifiec handler function
	   calling handleFunc add a rule to the multiplexer.
	   there is a default multiplexer, and you can write a custom multiplexer
	*/
}

// the *http.Request is, you guessed it, a pointer
// it is a pointer to the incoming http.request object...it is not a copy
func health_handler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("<h1>OKAY</h1>"))
    fmt.Println(request.Method)
}

/* handler function
    a ResponseWriter writes back to the client
    a Request gets intput from the client request
*/
func login_handler(writer http.ResponseWriter, request *http.Request) {
    fmt.Printf("login: %v\n", request.Method)
    writer.Header().Add("YO","mama")
    writer.WriteHeader(200)
}

func register_handler(writer http.ResponseWriter, request *http.Request){
    var user service_user
    err := json.NewDecoder(request.Body).Decode(&user)
    if err != nil{
        http.Error(writer, "invalid json", http.StatusBadRequest)
        return
    }
    if user.User_Name == "" || user.Password == "" {
        http.Error(writer, "username and password required", http.StatusBadRequest)
        return
    }

    hashed_pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
    if err != nil {
        http.Error(writer, "error hashing password", http.StatusInternalServerError)
        return
    }
    
    user.Password = string(hashed_pass)
    fmt.Println(request.RemoteAddr)
    ip, _, err := net.SplitHostPort(request.RemoteAddr)
    if err != nil{
        fmt.Printf("error with IP: %v\n",err)
        ip = request.RemoteAddr
    }
    user.IP_addr = ip
    user.Location = "Internet"
}

/*
    TO-DO
    register_hander first:
        request body. assign related values to a coherent struct
        check if uname exists, if so, then reject with good message. what error code is used here?
        if all clear, store to in-mem data store. 
    login_hander first:
        auth and return a success message
    
*/