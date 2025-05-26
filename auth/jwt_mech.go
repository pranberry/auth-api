package auth

import (
	"fmt"
	"jwt-auth/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
	secretKey,err := db.GetSecretKey()
	if err != nil{
		return JWTResponseStruct{}, err
	}
	tokenString, err := new_token.SignedString(secretKey)
	if err != nil {
		fmt.Printf("error generating JWT for %v: %v\n", username, err)
	}
	return JWTResponseStruct{ AccessToken: tokenString, TokenType: "bearer"}, err
}


func ValidateJWT(JWT string) (bool, error){

	claims := &jwt.RegisteredClaims{}
	secretKey, err := db.GetSecretKey()
	if err != nil{
		return false, err
	}
	keyFunc := func(token *jwt.Token) (any, error) {
		return secretKey, nil
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