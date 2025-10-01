package handlers

import (
	"auth-api/db"
	"auth-api/models"
	"encoding/json"
	"net"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	// registerGetUserByName and registerUserFunc allow tests to substitute
	// database access during handler execution.
	registerGetUserByName = db.GetUserByName
	registerUserFunc      = db.RegisterUser
)

// RegisterHandler handles POST /register requests and creates new user
// accounts when the payload is valid.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var resp = models.Response{
		Status:  http.StatusUnauthorized,
		Error:   nil,
		Message: "not allowed",
	}

	// close up request after return
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			resp.Error = err
			http.Error(w, resp.Message, resp.Status)
		}
	}()

	var user models.ServiceUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp.Error = err
		resp.Message = "invalid json"
		resp.Status = http.StatusBadRequest
		return
	}
	if user.Username == "" || user.Password == "" {
		resp.Message = "username and password required"
		resp.Status = http.StatusBadRequest
		return
	}

	// check if username exists in database
	// service user should be nil for an non-existent user
	serviceUser, err := registerGetUserByName(user.Username)
	if serviceUser != nil {
		resp.Message = "username taken. pick another"
		resp.Error = err
		resp.Status = http.StatusConflict
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.Message = "error hashing password"
		resp.Error = err
		resp.Status = http.StatusInternalServerError
		return
	}

	user.Password = string(hashedPass)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	user.IP_addr = ip
	user.Location = getLocation()

	err = registerUserFunc(user)
	if err != nil {
		resp.Error = err
		resp.Message = "failed to register user"
		resp.Status = http.StatusInternalServerError
		return
	} else {
		resp.Message = "user created successfully. proceed to login"
		resp.Status = http.StatusCreated
	}
}

func getLocation() string {
	return "Internet"
}
