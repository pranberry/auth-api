package user

import (
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)


/* handler function
	a http.Request gets intput from the client request...
    a http.ResponseWriter writes back to the client...headers, body, codes...
*/
func LoginHandler(writer http.ResponseWriter, request *http.Request) {
    var login_user_data ServiceUser
    err := json.NewDecoder(request.Body).Decode(&login_user_data)
    if err != nil {
        http.Error(writer, "Request Denied", http.StatusBadRequest)
        return        
    }

    user_data, user_exist := MasterUserDB[login_user_data.User_Name]
    if !user_exist{
        http.Error(writer, "username not found. register first", http.StatusBadRequest)
        return
    }

    // if user exists, validate password
    err = bcrypt.CompareHashAndPassword([]byte(user_data.Password), []byte(login_user_data.Password))
    if err != nil{
        http.Error(writer, "Password is incorrect", http.StatusBadRequest)
        return
    } else {
        resp_message := ResponseStruct{
            Message: "Login Successful",
            Username: user_data.User_Name,
        }
        writer.Header().Set("Content-type","application/json")
        writer.WriteHeader(http.StatusOK)
        json.NewEncoder(writer).Encode(resp_message)
    }
}