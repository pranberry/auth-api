package middleware

import (
	"auth-api/auth"
	"auth-api/handlers"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// statusRecorder helps us bring back the response's http.status-code
// Exclusively for middleware logging
// Embedding is awesome!
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// when WriteHeader is called, its updates both fields for statusRecorder
func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Log incoming Requests and outgoing responses
func Logger(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: %s, on %s, from %v", r.Method, r.URL.Path, r.RemoteAddr)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		log.Printf("RESPONSE: %d on %s, on %s, from %v", rec.status, r.Method, r.URL.Path, r.RemoteAddr)
	})
}

// Checks for an Authorization header and validates the token
func CheckJwt(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := handlers.Response{
			Status:  http.StatusUnauthorized,
			Message: "failed to validate auth token",
			Error:   nil,
		}

		// write response on fail
		defer handlers.WriteResponse(w, &resp)

		// retrieve header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			resp.Error = fmt.Errorf("no auth token")
			resp.Message = resp.Error.Error()
			resp.Status = http.StatusBadRequest
			return
		}
		authHeader, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			resp.Error = fmt.Errorf("corrupt token format")
			resp.Message = resp.Error.Error()
			resp.Status = http.StatusBadRequest
			return
		}
		err := auth.ValidateJWT(authHeader)
		if err != nil {
			resp.Error = fmt.Errorf("error validating token: %w", err)
			resp.Message = resp.Error.Error()
			return
		}

		next.ServeHTTP(w, r)
	})

}
