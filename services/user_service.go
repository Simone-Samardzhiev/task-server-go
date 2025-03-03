package services

import (
	"context"
	"net/http"
	"server/auth/passwords"
	"server/models"
	"server/repositories"
	"server/utils"
)

// UserService interface manage the business logic for users.
type UserService interface {
	// Register will check if everything is correct with the user information and return if not.
	Register(ctx context.Context, payload models.RegistrationsPayload) *utils.ErrorResponse
}

// DefaultUseService struct is the default implementation of [UserService].
type DefaultUseService struct {
	userRepository repositories.UserRepository
}

func (s *DefaultUseService) Register(ctx context.Context, payload models.RegistrationsPayload) *utils.ErrorResponse {
	result, err := s.userRepository.CheckIfEmailExists(ctx, payload.Email)
	if err != nil {
		return utils.InternalServerError()
	}

	if result {
		return utils.NewErrorResponse("Email already in use", http.StatusConflict)
	}

	result, err = s.userRepository.CheckIfUsernameExists(ctx, payload.Username)
	if err != nil {
		return utils.InternalServerError()
	}
	if result {
		return utils.NewErrorResponse("Username already in use", http.StatusConflict)
	}

	hash, err := passwords.HashPassword(payload.Password)
	if err != nil {
		return utils.InternalServerError()
	}

	err = s.userRepository.AddUser(ctx, payload.Email, payload.Username, hash)
	if err != nil {
		return utils.InternalServerError()
	}

	return nil
}

func NewDefaultService(userRepository repositories.UserRepository) *DefaultUseService {
	return &DefaultUseService{
		userRepository: userRepository,
	}
}
