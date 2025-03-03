package passwords

import (
	"testing"
)

// TestHashPassword will check if the [HashPassword] function works.
func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}

	t.Log(hash)
}

// TestVerifyPassword will check if the [VerifyPassword] function correctrly checks passwords.
func TestVerifyPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}

	t.Log(hash)

	result := VerifyPassword(password, hash)
	if !result {
		t.Errorf("VerifyPassword() doen't return true when checking the same password")
	}
}
