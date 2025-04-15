package models

import (
	"net/http"
	"server/utils"
	"strings"
)

// User struct holds user data.
type User struct {
	Id       int
	Email    string
	Username string
	Password string
}

// NewUser will create instance of [User]
func NewUser(id int, email, username, password string) *User {
	return &User{
		Id:       id,
		Email:    email,
		Username: username,
		Password: password,
	}
}

// LoginPayload is a struct holding login information.
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *LoginPayload) ValidatePayload() *utils.ErrorResponse {
	if p.Email == "" || p.Password == "" {
		return utils.NewErrorResponse("Invalid credentials", http.StatusUnauthorized)
	}

	return nil
}

// RegistrationsPayload is a struct holding the user information.
type RegistrationsPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *RegistrationsPayload) ValidatePayload() *utils.ErrorResponse {
	if res := u.ValidateEmail(); res != nil {
		return res
	}

	if res := u.ValidateEmail(); res != nil {
		return res
	}

	if res := u.ValidatePassword(); res != nil {
		return res
	}

	return nil
}

func (u *RegistrationsPayload) ValidateEmail() *utils.ErrorResponse {
	if strings.Contains(u.Email, " ") {
		return utils.NewErrorResponse("Email cannot contain spaces", http.StatusBadRequest)
	}

	if u.Email == "" {
		return utils.NewErrorResponse("Email is required", http.StatusBadRequest)
	}

	if !strings.Contains(u.Email, "@") {
		return utils.NewErrorResponse("Email does not contain @", http.StatusBadRequest)
	}

	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		return utils.NewErrorResponse("Email must contain @ only once", http.StatusBadRequest)
	}

	if parts[0] == "" {
		return utils.NewErrorResponse("Email is missing local part (before @)", http.StatusBadRequest)
	}

	if parts[1] == "" {
		return utils.NewErrorResponse("Email is missing remote part (after @)", http.StatusBadRequest)
	}

	if !strings.Contains(parts[1], ".") {
		return utils.NewErrorResponse("Domain must contain .", http.StatusBadRequest)
	}

	return nil
}

func (u *RegistrationsPayload) ValidateUsername() *utils.ErrorResponse {
	if strings.Contains(u.Username, " ") {
		return utils.NewErrorResponse("Username cannot contain spaces", http.StatusBadRequest)
	}

	if u.Username == "" {
		return utils.NewErrorResponse("Username is required", http.StatusBadRequest)
	}

	return nil
}

func (u *RegistrationsPayload) ValidatePassword() *utils.ErrorResponse {
	if strings.Contains(u.Password, " ") {
		return utils.NewErrorResponse("Password cannot contain spaces", http.StatusBadRequest)
	}

	if u.Password == "" {
		return utils.NewErrorResponse("Password is required", http.StatusBadRequest)
	}

	if len(u.Password) < 8 || len(u.Password) > 40 {
		return utils.NewErrorResponse("Password must be between 8 and 40 characters", http.StatusBadRequest)
	}

	if !strings.ContainsAny(u.Password, "-_@*&#!") {
		return utils.NewErrorResponse("Password must contain at least one special character", http.StatusBadRequest)
	}

	return nil
}
