package secret

import (
	"fmt"
	"jwt-auth/auth"
	"net/http"
	"strings"
)

// make a call to check the JWT token.
// then produce the "secret"
// which is an image of the hamster dance
func SecretHandler(writer http.ResponseWriter, request *http.Request) {

	// get the jwt from the auth header
	auth_header := request.Header.Get("Authorization")
	auth_header, found := strings.CutPrefix(auth_header,"Bearer ")
	if !found{
		http.Error(writer, "Corrupt Token Format", http.StatusUnauthorized)
		return
	}
	is_valid_token, err := auth.ValidateJWT(auth_header)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error validating token: %v", err), http.StatusUnauthorized)
		return
	}
	// don't really need this if block
	if is_valid_token {
		writer.Header().Set("Content-Type", "image/gif")
		//writer.WriteHeader(http.StatusOK)
		http.ServeFile(writer, request, "secret/hamster_dance.gif")
		return
	}
	
}