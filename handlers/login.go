package handlers

import (
	"encoding/json"
	"jwt-auth/auth"
	"jwt-auth/db"
	"jwt-auth/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	// loginGetUserByName and createJWTFunc are overridden in tests to avoid
	// hitting external dependencies.
	loginGetUserByName = db.GetUserByName
	createJWTFunc      = auth.CreateJWT
)

// LoginHandler processes POST /login requests and returns a JWT when the
// provided credentials are valid.
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var resp = models.Response{
		Status:  http.StatusUnauthorized,
		Error:   nil,
		Message: "not allowed",
	}
	
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			resp.Error = err
			http.Error(w, resp.Message, resp.Status)
		}
	}()

	var loginUserData models.ServiceUser
	err := json.NewDecoder(r.Body).Decode(&loginUserData)
	if err != nil {
		resp.Message = "request denied"
		resp.Status = http.StatusBadRequest
		resp.Error = err
		return
	}

	// check for user existence in db/mem
	userData, err := loginGetUserByName(loginUserData.Username)
	if err != nil {
		resp.Message = "username not found. register first"
		resp.Error = err
		resp.Status = http.StatusNotFound
		// return
	}

	// just in case check...
	if userData == nil{
		userData = &models.ServiceUser{
			Password: "123",
		}
	}
	
	// if user exists, validate password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginUserData.Password))
	if err != nil {
		resp.Message = "password is incorrect"
		resp.Error = err
		resp.Status = http.StatusBadRequest
		return
	}

	jwtResp, err := createJWTFunc(userData.Username)
	if err != nil {
		resp.Message = "failed to create jwt"
		resp.Status = http.StatusInternalServerError
		resp.Error = err
		return
	}
	resp.Message = "login successful"
	resp.Status = http.StatusOK
	resp.Error = nil
	resp.Data = jwtResp

}
