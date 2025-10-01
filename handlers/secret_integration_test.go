package handlers_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"auth-api/handlers"
	"auth-api/middleware"
)

func readBody(t *testing.T, res *http.Response) string {
	t.Helper()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return string(data)
}

func TestSecretHandlerLogsSuccess(t *testing.T) {
	var buf bytes.Buffer
	originalWriter := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(originalWriter)
	})

	handler := middleware.Logger(handlers.SecretHandler)

	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "image/gif" {
		t.Fatalf("expected gif content type, got %s", ct)
	}

	logs := buf.String()
	if !strings.Contains(logs, "RESPONSE: 200") {
		t.Fatalf("expected logs to contain status 200, got %q", logs)
	}
}

func TestServeStaticFileMissingAsset(t *testing.T) {
	var buf bytes.Buffer
	originalWriter := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(originalWriter)
	})

	handler := middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeStaticFile(w, r, "assets/not-real.gif", "image/gif")
	})

	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected json content type, got %s", ct)
	}
	body := readBody(t, res)
	if !strings.Contains(body, "requested asset not found") {
		t.Fatalf("unexpected body: %s", body)
	}

	logs := buf.String()
	if !strings.Contains(logs, "RESPONSE: 404") {
		t.Fatalf("expected logs to contain status 404, got %q", logs)
	}
}
