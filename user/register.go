package user

import (
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net"
)

/*
	- the *http.Request is, you guessed it, a pointer to the http.request object
*/
func RegisterHandler(writer http.ResponseWriter, request *http.Request){
    var user ServiceUser
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
    _, user_exist := MasterUserDB[user.User_Name]
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
    ip, _, err := net.SplitHostPort(request.RemoteAddr)
    if err != nil{
        ip = request.RemoteAddr
    }
    user.IP_addr = ip
    user.Location = "Internet"

    MasterUserDB[user.User_Name] = user
}