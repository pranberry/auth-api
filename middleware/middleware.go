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

		auth_header := r.Header.Get("Authorization")
		auth_header, ok := strings.CutPrefix(auth_header, "Bearer ")
		if !ok {
			log.Println("corrupt token format")
			//http.Error(writer, "corrupt token format", http.StatusUnauthorized)
		}
		err := auth.ValidateJWT(auth_header)
		if err != nil {
			err = fmt.Errorf("error validating token: %w\n", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	})

}
