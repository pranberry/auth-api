package user

import (
	"fmt"
	"jwt-auth/config"
	"jwt-auth/db"
	"jwt-auth/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Function takes a username and returns a valid JWT
func CreateJWT(username string) (models.ResponseStruct, error) {

	new_token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    config.TokenIssuer,
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JwtTTL)),
		})

	secretKey, err := db.GetSecretKey()
	if err != nil {
		return models.ResponseStruct{}, err
	}
	tokenString, err := new_token.SignedString(secretKey)
	if err != nil {
		fmt.Printf("error generating JWT for %v: %v\n", username, err)
	}
	return models.ResponseStruct{AccessToken: tokenString, TokenType: "bearer"}, err
}

// This function validates the JWT
func ValidateJWT(JWT string) (bool, error) {

	claims := &jwt.RegisteredClaims{}
	secretKey, err := db.GetSecretKey()
	if err != nil {
		return false, err
	}
	keyFunc := func(token *jwt.Token) (any, error) {
		return secretKey, nil
	}

	token, err := jwt.ParseWithClaims(JWT, claims, keyFunc)
	if err != nil {
		fmt.Printf("error: failed to parse token string: %v", err)
		return false, err
	}

	if token.Valid {
		// TODO: if token is valid, bump the expiry
		return true, nil
	} else {
		return false, fmt.Errorf("token is invalid")
	}
}
