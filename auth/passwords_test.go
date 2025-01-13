package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(&password)
	if err != nil {
		t.Error(err)
	}

	if hash == "" {
		t.Error("hash is empty")
	}

	if hash == password {
		t.Error("hash is the same")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	hash, err := HashPassword(&password)
	if err != nil {
		t.Error(err)
	}

	if !CheckPassword(&password, &hash) {
		t.Error("hash is different")
	}
}
