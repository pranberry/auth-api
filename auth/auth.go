package auth

import (
	"fmt"
	"jwt-auth/db"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// getSecretKey wraps db.GetSecretKey so tests can replace the dependency.
	getSecretKey = db.GetSecretKey
)

// JWTResponse represents the payload returned to clients after
// successfully authenticating.
type JWTResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Message     string `json:"message,omitempty"`
	HttpStatus  int
}

// CreateJWT creates a signed JWT for the provided username using the secret key
// stored in the database.
func CreateJWT(username string) (JWTResponse, error) {

	new_token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    getHostname(),
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(900 * time.Second)),
		})
	secretKey, err := getSecretKey()
	if err != nil {
		return JWTResponse{}, err
	}
	tokenString, err := new_token.SignedString(secretKey)
	if err != nil {
		fmt.Printf("error generating JWT for %v: %v\n", username, err)
	}
	return JWTResponse{AccessToken: tokenString, TokenType: "bearer"}, err
}

// ValidateJWT verifies the provided token string against the stored secret key.
func ValidateJWT(JWT string) error {

	claims := &jwt.RegisteredClaims{}
	secretKey, err := getSecretKey()
	if err != nil {
		return err
	}
	keyFunc := func(token *jwt.Token) (any, error) {
		return secretKey, nil
	}

	token, err := jwt.ParseWithClaims(JWT, claims, keyFunc)
	if err != nil {
		fmt.Printf("error: failed to parse token string: %v", err)
		return err
	}

	if token.Valid {
		// TODO: if token is valid, bump the expiry
		return nil
	} else {
		return fmt.Errorf("token is invalid")
	}
}

func getHostname() string {
	host, err := os.Hostname()
	if err != nil {
		host = "auth api"
	}
	return host
}
