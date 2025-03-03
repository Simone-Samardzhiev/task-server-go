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

// RegistrationsPayload is a struct holding the user information.
type RegistrationsPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *RegistrationsPayload) ValidatePayload() *utils.ValidateResponse {
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

func (u *RegistrationsPayload) ValidateEmail() *utils.ValidateResponse {
	if strings.Contains(u.Email, " ") {
		return &utils.ValidateResponse{
			Message:    "Email cannot contain spaces",
			StatusCode: http.StatusBadRequest,
		}
	}

	if u.Email == "" {
		return &utils.ValidateResponse{
			Message:    "Email is required",
			StatusCode: http.StatusBadRequest,
		}
	}

	if !strings.Contains(u.Email, "@") {
		return &utils.ValidateResponse{
			Message:    "Email must contain @",
			StatusCode: http.StatusBadRequest,
		}
	}

	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		return &utils.ValidateResponse{
			Message:    "Email must contain @ only once",
			StatusCode: http.StatusBadRequest,
		}
	}

	if parts[0] == "" {
		return &utils.ValidateResponse{
			Message:    "Email is missing local part (before @)",
			StatusCode: http.StatusBadRequest,
		}
	}

	if parts[1] == "" {
		return &utils.ValidateResponse{
			Message:    "Email is missing domain part (after @)",
			StatusCode: http.StatusBadRequest,
		}
	}

	if !strings.Contains(parts[1], ".") {
		return &utils.ValidateResponse{
			Message:    "Domain must contain .",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}

func (u *RegistrationsPayload) ValidateUsername() *utils.ValidateResponse {
	if strings.Contains(u.Username, " ") {
		return &utils.ValidateResponse{
			Message:    "Username cannot contain spaces",
			StatusCode: http.StatusBadRequest,
		}
	}

	if u.Username == "" {
		return &utils.ValidateResponse{
			Message:    "Username is required",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}

func (u *RegistrationsPayload) ValidatePassword() *utils.ValidateResponse {
	if strings.Contains(u.Password, " ") {
		return &utils.ValidateResponse{
			Message:    "Password cannot contain spaces",
			StatusCode: http.StatusBadRequest,
		}
	}

	if u.Password == "" {
		return &utils.ValidateResponse{
			Message:    "Password is required",
			StatusCode: http.StatusBadRequest,
		}
	}

	if len(u.Password) < 8 || len(u.Password) > 40 {
		return &utils.ValidateResponse{
			Message:    "Password must be between 8 and 40 characters",
			StatusCode: http.StatusBadRequest,
		}
	}

	if !strings.ContainsAny(u.Password, "-_@*&#!") {
		return &utils.ValidateResponse{
			Message:    "Password must contain at least one of this special characters(- _ @ * & # !)",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
