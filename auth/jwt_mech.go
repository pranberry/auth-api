package auth

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// secret key should live in the db...now that we have one
func SecretKey() []byte{
	var SecretKey string = "iexaiviazooJeiW0hex_o0O"
	return []byte(SecretKey)
}

type JWTResponseStruct struct{
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Message string `json:"message,omitempty"`
}

func CreateJWT(username string) (JWTResponseStruct, error) {
	new_token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer : "SCDP",
			Subject: username,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(900 * time.Second)),
		})
	tokenString, err := new_token.SignedString(SecretKey())
	if err != nil {
		fmt.Printf("error generating JWT for %v: %v\n", username, err)
	}
	return JWTResponseStruct{ AccessToken: tokenString, TokenType: "bearer"}, err
}


func ValidateJWT(JWT string) (bool, error){

	claims := &jwt.RegisteredClaims{}
	keyFunc := func(token *jwt.Token) (any, error) {
		return SecretKey(), nil
	}
	
	token, err := jwt.ParseWithClaims(JWT, claims, keyFunc)
	if err != nil{
		fmt.Printf("error: failed to parse token string: %v",err)
		return false, err
	}

	if token.Valid {
		// TODO: if token is valid, bump the expiry
		return true, nil
	}else{
		return false, fmt.Errorf("token is invalid")
	}
}