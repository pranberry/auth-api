package user

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
	default_test_user := ServiceUser {
		User_Name: legit_request_body_reusable["username"],
		Password: string(hashed_password),
		Location: "Internet",
		IP_addr: "127.0.0.1",
	}
	MasterUserDB[legit_request_body_reusable["username"]] = default_test_user
}

/*
	based on inputs, creates a http request and returns a http response object
	found myself writing the same code over and over agian...
*/
func request_response_helper(request_method string, endpoint string, request_body map[string]string) http.Response {

	json_body, _:= json.Marshal(request_body)
	request := httptest.NewRequest(request_method, endpoint, bytes.NewReader(json_body))
	writer := httptest.NewRecorder()

	switch endpoint {
	case "/register":
		RegisterHandler(writer, request)
	case "/login":
		LoginHandler(writer, request)
	}	
	
	response := writer.Result()

	return *response
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

	response := request_response_helper("POST", "/register",legit_request_body_reusable)

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

func Test_Register_New_User(t *testing.T){
	resetDB()

	new_user := map[string]string{
		"username": "billy_bob123",
		"password": "qwerty123",
	}
	response := request_response_helper("POST", "/register", new_user)

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, recieved: %d", response.StatusCode)
	}

	expected_message := "user created successfully"
	resp_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(resp_body_bytes)

	if !strings.Contains(resp_body_string, expected_message){
		t.Errorf("Expected message: %v, recieved: %v", expected_message, resp_body_string)
	}

}

func Test_Corrupt_Json_Body(t *testing.T){
	resetDB()

	junk_body := map[string]string{
		"user_name": "123449dkc",
		"pass_word": "fugazi",
	}
	response := request_response_helper("POST", "/register", junk_body)

	if response.StatusCode != http.StatusBadRequest{
		t.Errorf("Wrong status recieved for poorly formatted request. Expected 400. Recieved %d", response.StatusCode)
	}
}

func Test_Empty_Json_Body(t *testing.T){
	resetDB()

	new_body := map[string]string{
		"username": "",
		"password": "",
	}

	response := request_response_helper("POST", "/register", new_body)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected error code 400, recieved: %d", response.StatusCode)
	}

	expected_error := "username and password required"
	resp_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(resp_body_bytes)

	if !strings.Contains(resp_body_string, expected_error){
		t.Errorf("Expected message: %v, got: %v", expected_error, resp_body_string)
	}

}


//testing successful login
func Test_Login_Success(t *testing.T){

	resetDB()
	// just use the reuseable body to log-ing
	response := request_response_helper("POST", "/login", legit_request_body_reusable)

	if response.StatusCode != http.StatusOK{
		t.Errorf("Expected code 200, recieved: %d", response.StatusCode)
	}

	expected_resp := "Login Successful"
	resp_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(resp_body_bytes)

	if !strings.Contains(resp_body_string, expected_resp){
		t.Errorf("response message mismatch. expected: %v, recieved: %v", expected_resp, resp_body_string)
	}
}

// attempt login with non-existant user
func Test_user_non_exist(t *testing.T){
	resetDB()

	new_body := map[string]string{
		"username" : "bill_maher",
		"password" : "password123",
	}

	response := request_response_helper("POST", "/login", new_body)

	if response.StatusCode != http.StatusBadRequest{
		t.Errorf("Expected code 400, recieved: %d", response.StatusCode)
	}

	expected_message := "username not found. register first"
	resp_body_bytes, _ := io.ReadAll(response.Body)

	if !strings.Contains(string(resp_body_bytes),expected_message){
		t.Errorf("response message mismatch. expected: %v, recieved: %v", expected_message, string(resp_body_bytes))
	}
}

// testing incorrect password for existing username
func Test_incorrect_password(t *testing.T){
	resetDB()
	legit_request_body_reusable["password"] = "some_junk_pw123"

	response := request_response_helper("POST", "/login", legit_request_body_reusable)

	if response.StatusCode != http.StatusBadRequest{
		t.Errorf("Expected code 400, recieved: %d", response.StatusCode)
	}

	exp_msg := "Password is incorrect"
	resp_body_bytes, _ := io.ReadAll(response.Body)

	if !strings.Contains(string(resp_body_bytes), exp_msg){
		t.Errorf("response message mismatch. expected: %v, recieved: %v", exp_msg, string(resp_body_bytes))
	}

}