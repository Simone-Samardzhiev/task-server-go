package user

import "github.com/google/uuid"

// User is struct that defines used data.
type User struct {
	Id uuid.UUID
	// The email of the user.
	Email string
	// The password of the user.
	Password string
}

// JsonUser is a struct used to decode user information from requests.
type JsonUser struct {
	// The email of the user.
	Email string `json:"email"`
	// The password of the user.
	Password string `json:"password"`
}
