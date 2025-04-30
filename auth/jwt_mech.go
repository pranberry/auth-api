package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SecretKey() []byte{
	var SecretKey string = "iexaiviazooJeiW0hex_o0O"
	return []byte(SecretKey)
}

// I reckon i don't need this, can use the RegisterdClaims struct offered by JWT
type CustomClaim struct {
	Issuer		string 	`json:"iss"`
	Expiry		int64 	`json:"exp"`
	Issued_at	int64 	`json:"iat"`
	Username 	string 	`json:"username"`
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
	fmt.Println("RAW JWT: ", tokenString)
	if err != nil {
		fmt.Printf("error generating JWT for %v: %v\n", username, err)
	}
	return JWTResponseStruct{ AccessToken: tokenString, TokenType: "bearer"}, err
}


/*
	How to validate:
	- match signature
	- not expired yet
	- check if it was issued by you (iss == scdp)
	- is meant for username (sub==username)

*/
func ValidateJWT(JWT string) (bool, error){
	// read the auth header from the request
		// --- don't haave to do this with the jwt lib...
	// split it into its three parts
		//header.payload.signature
		// signing_input = header.payload
	// get signing mech from header
	// resign signing_input with secret key...
	// compare new sig with sig from request token
	// if shes good, then shes good

	claims := &jwt.RegisteredClaims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
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