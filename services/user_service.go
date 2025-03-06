package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"server/auth/passwords"
	"server/auth/tokens"
	"server/models"
	"server/repositories"
	"server/utils"
	"time"
)

// UserService interface manage the business logic for users.
type UserService interface {
	// Register will check if everything is correct with the user information and return if not.
	Register(ctx context.Context, payload models.RegistrationsPayload) *utils.ErrorResponse

	// Login will check used credentials and return group of token if user is authenticated.
	Login(ctx context.Context, payload models.LoginPayload) (*models.TokenGroup, *utils.ErrorResponse)
}

// DefaultUseService struct is the default implementation of [UserService].
type DefaultUseService struct {
	userRepository   repositories.UserRepository
	tokensRepository repositories.TokenRepository
	authenticator    *tokens.JWTAuthenticator
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

func (s *DefaultUseService) Login(ctx context.Context, payload models.LoginPayload) (*models.TokenGroup, *utils.ErrorResponse) {
	user, err := s.userRepository.GetUserByEmail(ctx, payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewErrorResponse("Invalid credentials", http.StatusUnauthorized)
	} else if err != nil {
		return nil, utils.InternalServerError()
	}

	passwordsMatch := passwords.VerifyPassword(payload.Password, user.Password)
	if !passwordsMatch {
		return nil, utils.NewErrorResponse("Invalid credentials", http.StatusUnauthorized)
	}

	tokenId := uuid.New()
	tokenExp := time.Now().Add(time.Hour * 24 * 7)
	refreshToken, err := s.authenticator.CreateRefreshToken(tokenId, tokenExp)
	if err != nil {
		return nil, utils.InternalServerError()
	}

	err = s.tokensRepository.AddToken(ctx, tokenId, tokenExp, user.Id)
	if err != nil {
		return nil, utils.InternalServerError()
	}

	accessToken, err := s.authenticator.CreateAccessToken(user.Id, time.Now().Add(time.Minute*10))
	if err != nil {
		return nil, utils.InternalServerError()
	}

	return models.NewTokenGroup(accessToken, refreshToken), nil
}

func NewDefaultService(userRepository repositories.UserRepository, tokenRepository repositories.TokenRepository, authenticator *tokens.JWTAuthenticator) *DefaultUseService {
	return &DefaultUseService{
		userRepository:   userRepository,
		tokensRepository: tokenRepository,
		authenticator:    authenticator,
	}
}
