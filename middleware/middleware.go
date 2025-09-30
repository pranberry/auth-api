package middleware

import (
	"fmt"
	"auth-api/auth"
	"log"
	"net/http"
	"strings"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: %s, on %s, from %v", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("RESPONSE: HTTPHEADERVALUEHERE on %s, on %s, from %v", r.Method, r.URL.Path, r.RemoteAddr)
		// how do i get the http.status sent from the request or write object?
	})
}

func CheckJwt(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no auth token", http.StatusUnauthorized)
			return
		}
		authHeader, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			log.Println("corrupt token format")
			http.Error(w, "corrupt token format", http.StatusUnauthorized)
			return
		}
		err := auth.ValidateJWT(authHeader)
		if err != nil {
			err = fmt.Errorf("error validating token: %w", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})

}
