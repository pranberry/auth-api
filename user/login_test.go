package user

import (
	"bytes"
	"encoding/json"
	"io"
	"jwt-auth/config"
	"jwt-auth/db"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var legit_request_body_reusable = map[string]string{
	"username": "joe_smith_TU",
	"password": "test_pass_123",
}

/*
based on inputs, creates a http request and returns a http response object
found myself writing the same code over and over agian...
*/
func request_response_helper(request_method string, endpoint string, request_body map[string]string) http.Response {

	json_body, _ := json.Marshal(request_body)
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
NOTE:
- TestMain is a special function in the testing package. it tests the test harness itself
- Notice the (m *testing.M) in the sig...compared to the (t testing.T) in the regular test cases
- the M struct is for the Test Manager, and the T struct is individual test context
*/
func TestMain(m *testing.M) {
	println("test main running")
	err := db.InitDB(config.User, config.TestDb, config.Password, config.Host)
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

// testing successful login
func Test_Login_Success(t *testing.T) {

	// just use the reuseable body to log-ing
	response := request_response_helper("POST", "/login", legit_request_body_reusable)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected code 200, recieved: %d", response.StatusCode)
	}

	expected_resp := "login successful"
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

	response := request_response_helper("POST", "/login", new_body)

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

	legit_request_body_reusable["password"] = "some_junk_pw123"

	response := request_response_helper("POST", "/login", legit_request_body_reusable)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected code 400, recieved: %d", response.StatusCode)
	}

	exp_msg := "incorrect password"
	resp_body_bytes, _ := io.ReadAll(response.Body)

	if !strings.Contains(string(resp_body_bytes), exp_msg) {
		t.Errorf("response message mismatch. expected: %v, recieved: %v", exp_msg, string(resp_body_bytes))
	}

}
