package user

import (
	"encoding/json"
	"jwt-auth/db"
	"jwt-auth/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Handle Logic for Login attempt. Returns a JWT on successful login
func LoginHandler(writer http.ResponseWriter, request *http.Request) {

	var login_user_data models.ServiceUser
	err := json.NewDecoder(request.Body).Decode(&login_user_data)
	if err != nil {
		SendReponse(models.ResponseStruct{
			Message: "flawed request",
		},
			http.StatusBadRequest,
			writer,
		)
		return
	}

	// check for user existance in db
	user_data, err := db.GetUserByName(login_user_data.User_Name)
	if err != nil {
		SendReponse(models.ResponseStruct{
			Message: "username not found. register first",
		},
			http.StatusBadRequest,
			writer)
		return
	}

	// if user exists, validate password
	err = bcrypt.CompareHashAndPassword([]byte(user_data.Password), []byte(login_user_data.Password))
	if err != nil {
		SendReponse(models.ResponseStruct{
			Message: "incorrect password",
		},
			http.StatusBadRequest,
			writer)
		return
	} else {
		jwt_resp, err := CreateJWT(user_data.User_Name)
		if err != nil {
			SendReponse(models.ResponseStruct{
				Message: "error createing JWT: " + err.Error(),
			},
				http.StatusInternalServerError,
				writer)
			return
		}
		jwt_resp.Message = "login successful"
		SendReponse(jwt_resp, http.StatusOK, writer)
		return
	}
}
