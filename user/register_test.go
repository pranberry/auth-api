package user

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_Register_Existing_User(t *testing.T) {

	response := request_response_helper("POST", "/register", legit_request_body_reusable)

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
	response := request_response_helper("POST", "/register", new_user)

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
	response := request_response_helper("POST", "/register", junk_body)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status recieved for poorly formatted request. Expected 400. Recieved %d", response.StatusCode)
	}
}

func Test_Empty_Json_Body(t *testing.T) {

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

	if !strings.Contains(resp_body_string, expected_error) {
		t.Errorf("Expected message: %v, got: %v", expected_error, resp_body_string)
	}

}
