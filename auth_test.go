package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)


var legit_request_body_reusable = map[string]string{
	"username": "joe_smith_TU",
	"password": "test_pass_123",
}

func resetDB(){
	hashed_password, _ := bcrypt.GenerateFromPassword([]byte(legit_request_body_reusable["password"]), bcrypt.DefaultCost)
	default_test_user := service_user {
		User_Name: "joe_smith_TU",
		Password: string(hashed_password),
		Location: "Internet",
		IP_addr: "127.0.0.1",
	}
	master_user_db[default_test_user.User_Name] = default_test_user
}

/* 
	Things to test:
	- testing with empty/incomplete json-body
	- testing with mumbo-jumbo as json-body
	- testing with existing username
	- testing successful register
*/
func Test_Register_Existing_User(t *testing.T){
	resetDB()
	json_body, _ := json.Marshal(legit_request_body_reusable)

	request := httptest.NewRequest("POST","/register", bytes.NewReader(json_body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	register_handler(writer, request)
	response := writer.Result()

	if response.StatusCode != http.StatusBadRequest{
		t.Errorf("Expected status 400, got: %v", response.StatusCode)
	}

	expected_message := "username taken. pick another"
	response_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(response_body_bytes)

	if !strings.Contains(resp_body_string, expected_message){
		t.Errorf("expected response body to contain: %v, got: %v", expected_message, resp_body_string)
	}

}


func TestLogin(t *testing.T){
	/*
	Things to Test:
	- testing junk as json-body
	- testing non-existing username
	- testing incorrect password for existing username
	- testing successful login
	*/

}


