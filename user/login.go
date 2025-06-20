package user

import (
	"encoding/json"
	"fmt"
	"jwt-auth/auth"
	"jwt-auth/db"
	"jwt-auth/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

/*
	 handler function
		a http.Request gets intput from the client request...
	    a http.ResponseWriter writes back to the client...headers, body, codes...
*/
func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	var login_user_data models.ServiceUser
	err := json.NewDecoder(request.Body).Decode(&login_user_data)
	if err != nil {
		http.Error(writer, "Request Denied", http.StatusBadRequest)
		return
	}
	// check for user existance in db/mem
	user_data, err := db.GetUserByName(login_user_data.User_Name)

	if err != nil {
		http.Error(writer, "username not found. register first", http.StatusBadRequest)
		return
	}

	// if user exists, validate password
	err = bcrypt.CompareHashAndPassword([]byte(user_data.Password), []byte(login_user_data.Password))
	if err != nil {
		http.Error(writer, "Password is incorrect", http.StatusBadRequest)
		return
	} else {
		jwt_resp, err := auth.CreateJWT(user_data.User_Name)
		if err != nil {
			err := fmt.Sprintf("error creating jwt: %v", err)
			http.Error(writer, err, http.StatusInternalServerError)
			return
		}
		jwt_resp.Message = "Login Successful"
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(jwt_resp)
		return

	}
}
