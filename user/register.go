package user

import (
	"encoding/json"
	"jwt-auth/db"
	"jwt-auth/models"
	"net"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

/*
	- the *http.Request is, you guessed it, a pointer to the http.request object
*/
func RegisterHandler(writer http.ResponseWriter, request *http.Request){
    var user models.ServiceUser
    err := json.NewDecoder(request.Body).Decode(&user)
    if err != nil{
        http.Error(writer, "invalid json", http.StatusBadRequest)
        return
    }
    if user.User_Name == "" || user.Password == "" {
        http.Error(writer, "username and password required", http.StatusBadRequest)
        return
    }

    // check if username exists in database
    service_user, _ := db.GetUserByName(user.User_Name)

    if service_user != nil {// service user should be nil for an non-existiant user 
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


    err = db.RegisterUser(user)
    if err != nil{
        http.Error(writer, "Failed to register user", http.StatusInternalServerError)
        return
    }else{
        resp := models.ResponseStruct{}
        resp.Message = "user created successfully"
        resp.Username = user.User_Name
        writer.Header().Set("Content-Type","application/json")
        writer.WriteHeader(http.StatusCreated)
        json.NewEncoder(writer).Encode(resp)
    }
}