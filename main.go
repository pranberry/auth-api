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
var master_user_db = make(map[string]service_user)

type ResponseStruct struct {
    Message string `json:"message"`
    Username string `json:"username,omitempty"`
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
	writer.Write([]byte("<h1>ALIVE AND WELL...ish</h1>"))
    fmt.Println(request.Method)
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

    // check if user_name
    _, user_exist := master_user_db[user.User_Name]
    if user_exist {
        http.Error(writer, "username taken. pick another", http.StatusBadRequest)
        return
    }

    hashed_pass, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
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

    master_user_db[user.User_Name] = user
}

/* handler function
    a ResponseWriter writes back to the client
    a Request gets intput from the client request
*/
func login_handler(writer http.ResponseWriter, request *http.Request) {
    var login_user_data service_user
    err := json.NewDecoder(request.Body).Decode(&login_user_data)
    if err != nil {
        http.Error(writer, "Request Denied", http.StatusBadRequest)
        return        
    }

    user_data, user_exist := master_user_db[login_user_data.User_Name]
    if !user_exist{
        http.Error(writer, "username not found. register first", http.StatusBadRequest)
        return
    }

    // if user exists, validate password
    err = bcrypt.CompareHashAndPassword([]byte(user_data.Password), []byte(login_user_data.Password))
    if err != nil{
        http.Error(writer, "Password is incorrect", http.StatusBadRequest)
        fmt.Printf("error: %v\n", err)
        return
    } else {
        fmt.Println("Password and hash match...login successful")
        resp_message := ResponseStruct{
            Message: "Login Successful",
            Username: user_data.User_Name,
        }
        writer.Header().Set("Content-type","application/json")
        writer.WriteHeader(http.StatusOK)
        json.NewEncoder(writer).Encode(resp_message)
    }



}

/*
    TO-DO
    login_hander first:
        auth and return a success message
    - write tests
    - change all error messages to send back JSON, not plaintext
    - eventually, add https to the mix, so we're getting encrypted data
*/