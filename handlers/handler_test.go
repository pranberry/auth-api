package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jwt-auth/db"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var legitRequestBodyReusable = map[string]string{
	"username": "joe_smith_TU",
	"password": "test_pass_123",
}

/*
based on inputs, creates a http request and returns a http response object
found myself writing the same code over and over agian...
format:
	requestResponseHelper("POST", "/register", legitRequestBodyReusable)
*/
func requestResponseHelper(rMethod string, url string, request_body map[string]string) http.Response {

	json_body, _ := json.Marshal(request_body)
	request := httptest.NewRequest(rMethod, url, bytes.NewReader(json_body))
	writer := httptest.NewRecorder()

	switch url {
	case "/register":
		RegisterHandler(writer, request)
	case "/login":
		LoginHandler(writer, request)
	}

	response := writer.Result()

	return *response
}

/*
NOTE:
- TestMain is a special function in the testing package. it tests the test harness itself
- Notice the (m *testing.M) in the sig...compared to the (t testing.T) in the regular test cases
- the M struct is for the Test Manager, and the T struct is individual test context
*/
func TestMain(m *testing.M) {
	err := db.InitDB("token_master", "jwt_test", "tokenPass", "auth-db")
	if err != nil {
		log.Fatal("DB init failed: ", err)
	}
	os.Exit(m.Run())
}


func Test_Register_Existing_User(t *testing.T) {

	response := requestResponseHelper("POST", "/register", legitRequestBodyReusable)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got: %v", response.StatusCode)
	}

	expected_message := "username taken. pick another"
	response_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(response_body_bytes)

	if !strings.Contains(resp_body_string, expected_message) {
		t.Errorf("expected response body to contain: %v, got: %v", expected_message, resp_body_string)
	}

}

func Test_Register_New_User(t *testing.T) {

	// since the change to DB, the username, once created, will persist between the tests.
	// append some random string, a count of nanaoseconds to the username. unique everytime.
	username := fmt.Sprintf("billy_bob123_%v", time.Now().UnixNano())
	new_user := map[string]string{
		"username": username,
		"password": "qwerty123",
	}
	response := requestResponseHelper("POST", "/register", new_user)

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, recieved: %d", response.StatusCode)
	}

	expected_message := "user created successfully"
	resp_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(resp_body_bytes)

	if !strings.Contains(resp_body_string, expected_message) {
		t.Errorf("Expected message: %v, recieved: %v", expected_message, resp_body_string)
	}

}

func Test_Corrupt_Json_Body(t *testing.T) {

	junk_body := map[string]string{
		"user_name": "123449dkc",
		"pass_word": "fugazi",
	}
	response := requestResponseHelper("POST", "/register", junk_body)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status recieved for poorly formatted request. Expected 400. Recieved %d", response.StatusCode)
	}
}

func Test_Empty_Json_Body(t *testing.T) {

	new_body := map[string]string{
		"username": "",
		"password": "",
	}

	response := requestResponseHelper("POST", "/register", new_body)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected error code 400, recieved: %d", response.StatusCode)
	}

	expected_error := "username and password required"
	resp_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(resp_body_bytes)

	if !strings.Contains(resp_body_string, expected_error) {
		t.Errorf("Expected message: %v, got: %v", expected_error, resp_body_string)
	}

}

// testing successful login
func Test_Login_Success(t *testing.T) {

	// just use the reuseable body to log-ing
	response := requestResponseHelper("POST", "/login", legitRequestBodyReusable)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected code 200, recieved: %d", response.StatusCode)
	}

	expected_resp := "Login Successful"
	resp_body_bytes, _ := io.ReadAll(response.Body)
	resp_body_string := string(resp_body_bytes)

	if !strings.Contains(resp_body_string, expected_resp) {
		t.Errorf("response message mismatch. expected: %v, recieved: %v", expected_resp, resp_body_string)
	}
}

// attempt login with non-existant user
func Test_user_non_exist(t *testing.T) {

	new_body := map[string]string{
		"username": "bill_maher",
		"password": "password123",
	}

	response := requestResponseHelper("POST", "/login", new_body)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected code 400, recieved: %d", response.StatusCode)
	}

	expected_message := "username not found. register first"
	resp_body_bytes, _ := io.ReadAll(response.Body)

	if !strings.Contains(string(resp_body_bytes), expected_message) {
		t.Errorf("response message mismatch. expected: %v, recieved: %v", expected_message, string(resp_body_bytes))
	}
}

// testing incorrect password for existing username
func Test_incorrect_password(t *testing.T) {

	legitRequestBodyReusable["password"] = "some_junk_pw123"

	response := requestResponseHelper("POST", "/login", legitRequestBodyReusable)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected code 400, recieved: %d", response.StatusCode)
	}

	exp_msg := "Password is incorrect"
	resp_body_bytes, _ := io.ReadAll(response.Body)

	if !strings.Contains(string(resp_body_bytes), exp_msg) {
		t.Errorf("response message mismatch. expected: %v, recieved: %v", exp_msg, string(resp_body_bytes))
	}

}
