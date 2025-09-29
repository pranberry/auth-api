package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestCreateJWTSuccess ensures that CreateJWT returns a signed token and the
// expected metadata when the secret key lookup succeeds.
func TestCreateJWTSuccess(t *testing.T) {
	originalGetSecretKey := getSecretKey
	getSecretKey = func() ([]byte, error) {
		return []byte("secret"), nil
	}
	t.Cleanup(func() {
		getSecretKey = originalGetSecretKey
	})

	// Generate a token for a known user and validate the response contract.
	resp, err := CreateJWT("alice")
	if err != nil {
		t.Fatalf("CreateJWT returned unexpected error: %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatalf("expected access token to be populated")
	}
	if resp.TokenType != "bearer" {
		t.Fatalf("expected token type 'bearer', got %s", resp.TokenType)
	}

	// Independently parse the JWT to ensure the claims were encoded correctly.
	token, err := jwt.ParseWithClaims(resp.AccessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse generated token: %v", err)
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		t.Fatalf("unexpected claim type: %T", token.Claims)
	}
	if claims.Subject != "alice" {
		t.Errorf("expected subject 'alice', got %s", claims.Subject)
	}
	if claims.Issuer != "SCDP" {
		t.Errorf("expected issuer 'SCDP', got %s", claims.Issuer)
	}
	if claims.ExpiresAt == nil || time.Until(claims.ExpiresAt.Time) > 16*time.Minute || time.Until(claims.ExpiresAt.Time) < 14*time.Minute {
		t.Errorf("expected expiry about 15 minutes from now, got %v", claims.ExpiresAt)
	}
}

// TestCreateJWTSecretKeyError confirms that CreateJWT propagates failures when
// the signing secret cannot be retrieved.
func TestCreateJWTSecretKeyError(t *testing.T) {
	originalGetSecretKey := getSecretKey
	getSecretKey = func() ([]byte, error) {
		return nil, errors.New("boom")
	}
	t.Cleanup(func() {
		getSecretKey = originalGetSecretKey
	})

	_, err := CreateJWT("alice")
	if err == nil {
		t.Fatalf("expected error when secret key retrieval fails")
	}
}

// TestValidateJWT validates that well-formed tokens pass verification while
// malformed tokens fail.
func TestValidateJWT(t *testing.T) {
	originalGetSecretKey := getSecretKey
	getSecretKey = func() ([]byte, error) {
		return []byte("secret"), nil
	}
	t.Cleanup(func() {
		getSecretKey = originalGetSecretKey
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:  "alice",
		Issuer:   "SCDP",
		IssuedAt: jwt.NewNumericDate(time.Now()),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	if err := ValidateJWT(tokenString); err != nil {
		t.Fatalf("expected token to be valid, got error: %v", err)
	}

	if err := ValidateJWT("not-a-token"); err == nil {
		t.Fatalf("expected validation error for malformed token")
	}
}

// TestValidateJWTSecretKeyError ensures an error from the secret key lookup is
// returned to the caller.
func TestValidateJWTSecretKeyError(t *testing.T) {
	originalGetSecretKey := getSecretKey
	getSecretKey = func() ([]byte, error) {
		return nil, errors.New("boom")
	}
	t.Cleanup(func() {
		getSecretKey = originalGetSecretKey
	})

	if err := ValidateJWT("anything"); err == nil {
		t.Fatalf("expected error when secret key lookup fails")
	}
}
