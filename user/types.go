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

func NewUser(email string, id *uuid.UUID, password string) *User {
	return &User{Email: email, Id: *id, Password: password}
}

// JsonUser is a struct used to decode user information from requests.
type JsonUser struct {
	// The email of the user.
	Email string `json:"email"`
	// The password of the user.
	Password string `json:"password"`
}

func NewJsonUser(email string, password string) *JsonUser {
	return &JsonUser{
		Email:    email,
		Password: password,
	}
}
