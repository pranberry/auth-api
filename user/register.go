package user

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"jwt-auth/db"
	"jwt-auth/models"
	"net"
	"net/http"
)

/*
Handler function for registering a user
*/
func RegisterHandler(writer http.ResponseWriter, request *http.Request) {

	var user models.ServiceUser
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, "invalid json", http.StatusBadRequest)
		return
	}
	if user.User_Name == "" || user.Password == "" {
		http.Error(writer, "username and password required", http.StatusBadRequest)
		return
	}

	// check if username exists in database
	userExists, err := db.CheckUserExists(user.User_Name)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// If the users exists, return error
	if userExists {
		SendReponse(models.ResponseStruct{
			Message: "username taken. pick another",
		},
			http.StatusBadRequest,
			writer)

		return
	}

	// hash the password from the request
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		SendReponse(models.ResponseStruct{Message: "error hashing password"}, http.StatusInternalServerError, writer)
		return
	}
	// Replace string password with bycrypt hash
	user.Password = string(hashed_pass)

	// Get IP and location value...this is largely useless
	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		ip = request.RemoteAddr
	}
	user.IP_addr = ip
	user.Location = "Internet"

	// Register the user
	err = db.RegisterUser(user)
	if err != nil {
		SendReponse(
			models.ResponseStruct{
				Message: "failed to register user: " + err.Error(),
			},
			http.StatusInternalServerError,
			writer,
		)
		return
	} else {
		resp := models.ResponseStruct{
			Message: "user created successfully",
			Username: user.User_Name,
		}
		SendReponse(resp, http.StatusCreated, writer)
		return
	}
}

// Helper to format the response sent back to the user
func SendReponse(resp models.ResponseStruct, httpRespCode int, writer http.ResponseWriter) {

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpRespCode)
	json.NewEncoder(writer).Encode(resp)

}