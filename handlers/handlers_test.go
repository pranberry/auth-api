package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-api/auth"
	"auth-api/models"

	"golang.org/x/crypto/bcrypt"
)

// newJSONRequest builds a test HTTP request with a JSON payload and consistent
// remote address metadata.
func newJSONRequest(t *testing.T, method, target string, body any) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, target, &buf)
	req.RemoteAddr = "127.0.0.1:12345"
	return req
}

// readBody provides a convenience wrapper to read response bodies in tests.
func readBody(t *testing.T, res *http.Response) string {
	t.Helper()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return string(data)
}

// TestRegisterHandlerSuccess checks that a valid registration request succeeds
// and persists the transformed user payload.
func TestRegisterHandlerSuccess(t *testing.T) {
	originalGet := registerGetUserByName
	originalRegister := registerUserFunc
	registerGetUserByName = func(username string) (*models.ServiceUser, error) {
		return nil, errors.New("not found")
	}
	registerUserFunc = func(user models.ServiceUser) error {
		if user.Username != "alice" || user.Location != "Internet" || user.IP_addr == "" {
			t.Fatalf("unexpected user payload: %+v", user)
		}
		return nil
	}
	t.Cleanup(func() {
		registerGetUserByName = originalGet
		registerUserFunc = originalRegister
	})

	req := newJSONRequest(t, http.MethodPost, "/register", map[string]string{
		"username": "alice",
		"password": "password123",
	})
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	res := rr.Result()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content type: %s", ct)
	}
	body := readBody(t, res)
	if !bytes.Contains([]byte(body), []byte("user created successfully")) {
		t.Fatalf("expected success message, got %s", body)
	}
}

// TestRegisterHandlerUsernameTaken ensures an existing username results in a
// user-friendly error.
func TestRegisterHandlerUsernameTaken(t *testing.T) {
	originalGet := registerGetUserByName
	registerGetUserByName = func(username string) (*models.ServiceUser, error) {
		return &models.ServiceUser{Username: username}, nil
	}
	t.Cleanup(func() {
		registerGetUserByName = originalGet
	})

	req := newJSONRequest(t, http.MethodPost, "/register", map[string]string{
		"username": "alice",
		"password": "password123",
	})
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	res := rr.Result()
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
	body := readBody(t, res)
	if !bytes.Contains([]byte(body), []byte("username taken")) {
		t.Fatalf("unexpected body: %s", body)
	}
}

// TestRegisterHandlerInvalidJSON verifies malformed JSON short-circuits
// request processing.
func TestRegisterHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", rr.Result().StatusCode)
	}
}

// TestRegisterHandlerMissingFields confirms the handler validates required
// fields before touching persistence.
func TestRegisterHandlerMissingFields(t *testing.T) {
	req := newJSONRequest(t, http.MethodPost, "/register", map[string]string{})
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing fields, got %d", rr.Result().StatusCode)
	}
}

// TestRegisterHandlerPersistenceError makes sure database failures are
// translated to a 500 response.
func TestRegisterHandlerPersistenceError(t *testing.T) {
	originalGet := registerGetUserByName
	originalRegister := registerUserFunc
	registerGetUserByName = func(username string) (*models.ServiceUser, error) {
		return nil, errors.New("not found")
	}
	registerUserFunc = func(user models.ServiceUser) error {
		return errors.New("db down")
	}
	t.Cleanup(func() {
		registerGetUserByName = originalGet
		registerUserFunc = originalRegister
	})

	req := newJSONRequest(t, http.MethodPost, "/register", map[string]string{
		"username": "alice",
		"password": "password123",
	})
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 for persistence error, got %d", rr.Result().StatusCode)
	}
}

// TestLoginHandlerSuccess covers the happy path including password validation
// and JWT issuance.
func TestLoginHandlerSuccess(t *testing.T) {
	originalGet := loginGetUserByName
	originalCreate := createJWTFunc
	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	loginGetUserByName = func(username string) (*models.ServiceUser, error) {
		return &models.ServiceUser{Username: username, Password: string(hashed)}, nil
	}
	createJWTFunc = func(username string) (auth.JWTResponse, error) {
		return auth.JWTResponse{AccessToken: "token", TokenType: "bearer"}, nil
	}
	t.Cleanup(func() {
		loginGetUserByName = originalGet
		createJWTFunc = originalCreate
	})

	req := newJSONRequest(t, http.MethodPost, "/login", map[string]string{
		"username": "alice",
		"password": "password123",
	})
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	res := rr.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	body := readBody(t, res)
	if !bytes.Contains([]byte(body), []byte("Login Successful")) {
		t.Fatalf("unexpected body: %s", body)
	}
}

// TestLoginHandlerInvalidJSON validates malformed JSON is rejected.
func TestLoginHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", rr.Result().StatusCode)
	}
}

// TestLoginHandlerUserNotFound ensures missing users result in a clear
// bad-request error.
func TestLoginHandlerUserNotFound(t *testing.T) {
	originalGet := loginGetUserByName
	loginGetUserByName = func(username string) (*models.ServiceUser, error) {
		return nil, errors.New("not found")
	}
	t.Cleanup(func() {
		loginGetUserByName = originalGet
	})

	req := newJSONRequest(t, http.MethodPost, "/login", map[string]string{
		"username": "alice",
		"password": "password123",
	})
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing user, got %d", rr.Result().StatusCode)
	}
}

// TestLoginHandlerInvalidPassword checks that incorrect credentials are
// rejected with a 400 response.
func TestLoginHandlerInvalidPassword(t *testing.T) {
	originalGet := loginGetUserByName
	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	loginGetUserByName = func(username string) (*models.ServiceUser, error) {
		return &models.ServiceUser{Username: username, Password: string(hashed)}, nil
	}
	t.Cleanup(func() {
		loginGetUserByName = originalGet
	})

	req := newJSONRequest(t, http.MethodPost, "/login", map[string]string{
		"username": "alice",
		"password": "wrong",
	})
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad password, got %d", rr.Result().StatusCode)
	}
}

// TestLoginHandlerTokenFailure verifies JWT creation errors surface as a 500
// response.
func TestLoginHandlerTokenFailure(t *testing.T) {
	originalGet := loginGetUserByName
	originalCreate := createJWTFunc
	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	loginGetUserByName = func(username string) (*models.ServiceUser, error) {
		return &models.ServiceUser{Username: username, Password: string(hashed)}, nil
	}
	createJWTFunc = func(username string) (auth.JWTResponse, error) {
		return auth.JWTResponse{}, errors.New("fail")
	}
	t.Cleanup(func() {
		loginGetUserByName = originalGet
		createJWTFunc = originalCreate
	})

	req := newJSONRequest(t, http.MethodPost, "/login", map[string]string{
		"username": "alice",
		"password": "password123",
	})
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 when token creation fails, got %d", rr.Result().StatusCode)
	}
}

// TestSecretHandler ensures the static secret endpoint returns the expected
// status and content type.
func TestSecretHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	rr := httptest.NewRecorder()

	SecretHandler(rr, req)

	res := rr.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "image/gif" {
		t.Fatalf("expected gif content type, got %s", ct)
	}
}
