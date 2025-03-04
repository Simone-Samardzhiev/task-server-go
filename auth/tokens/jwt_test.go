package tokens

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

var authenticator = NewJWTAuthenticator([]byte("secret"), "issuer")

func TestJWTAuthenticatorCreateRefreshToken(t *testing.T) {
	token, err := authenticator.CreateRefreshToken(uuid.New(), time.Now().Add(time.Hour*24*14))
	if err != nil {
		t.Fatalf("Error creating refresh token: %v", err)
	}

	t.Logf("Got refresh token token: %v", token)
}

func TestJWTAuthenticatorCreateAccessToken(t *testing.T) {
	token, err := authenticator.CreateAccessToken(1, time.Now().Add(time.Minute*10))
	if err != nil {
		t.Fatalf("Error creating access token token: %v", err)
	}
	t.Logf("Got access token token: %v", token)
}

func TestJWTAuthenticatorVerifyRefreshToken(t *testing.T) {
	// Create a new token.
	token, err := authenticator.CreateRefreshToken(uuid.New(), time.Now().Add(time.Minute*10))
	if err != nil {
		t.Fatalf("Error creating refresh token: %v", err)
	}
	t.Logf("Got refresh token token: %v", token)

	// Check if verify works.
	claims, err := authenticator.VerifyToken(token, RefreshTokenType)
	if err != nil {
		t.Fatalf("Error verifying access token: %v", err)
	}
	t.Logf("Got refresh token claims: %v", claims)

	// Check if verify return error if the types doesn't match.
	claims, err = authenticator.VerifyToken(token, AccessTokenType)
	if err == nil {
		t.Fatal("Expected error, because the type for verification is wrong")
	}
}

func TestJWTAuthenticatorVerifyAccessToken(t *testing.T) {
	// Create a new token.
	token, err := authenticator.CreateAccessToken(1, time.Now().Add(time.Minute*10))
	if err != nil {
		t.Fatalf("Error creating access token: %v", err)
	}
	t.Logf("Got access token token: %v", token)

	// Check if verify works.
	claims, err := authenticator.VerifyToken(token, AccessTokenType)
	if err != nil {
		t.Fatalf("Error verifying access token: %v", err)
	}
	t.Logf("Got access token claims: %v", claims)

	// Check if verify return error if the types doesn't match.
	claims, err = authenticator.VerifyToken(token, RefreshTokenType)
	if err == nil {
		t.Fatal("Expected error, because the type for verification is wrong")
	}
}
