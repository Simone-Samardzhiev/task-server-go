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

	// Refresh will check if the token is valid. If the token is valid
	// it will be deleted and new refresh token and access token will be generated.
	Refresh(ctx context.Context, token tokens.Token) (*models.TokenGroup, *utils.ErrorResponse)
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
		return utils.InternalServerErrorResponse()
	}

	if result {
		return utils.NewErrorResponse("Email already in use", http.StatusConflict)
	}

	result, err = s.userRepository.CheckIfUsernameExists(ctx, payload.Username)
	if err != nil {
		return utils.InternalServerErrorResponse()
	}
	if result {
		return utils.NewErrorResponse("Username already in use", http.StatusConflict)
	}

	hash, err := passwords.HashPassword(payload.Password)
	if err != nil {
		return utils.InternalServerErrorResponse()
	}

	err = s.userRepository.AddUser(ctx, payload.Email, payload.Username, hash)
	if err != nil {
		return utils.InternalServerErrorResponse()
	}

	return nil
}

// createTokenGroup will create a refresh token add it to the database and create an access token.
func (s *DefaultUseService) createTokenGroup(ctx context.Context, userId int) (*models.TokenGroup, *utils.ErrorResponse) {
	tokenId := uuid.New()
	tokenExp := time.Now().Add(time.Hour * 24 * 7)
	refreshToken, err := s.authenticator.CreateRefreshToken(tokenId, tokenExp)
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	err = s.tokensRepository.AddToken(ctx, tokenId, tokenExp, userId)
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	accessToken, err := s.authenticator.CreateAccessToken(userId, time.Now().Add(time.Minute*10))
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	return models.NewTokenGroup(accessToken, refreshToken), nil
}

func (s *DefaultUseService) Login(ctx context.Context, payload models.LoginPayload) (*models.TokenGroup, *utils.ErrorResponse) {
	user, err := s.userRepository.GetUserByEmail(ctx, payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewErrorResponse("Invalid credentials", http.StatusUnauthorized)
	} else if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	passwordsMatch := passwords.VerifyPassword(payload.Password, user.Password)
	if !passwordsMatch {
		return nil, utils.NewErrorResponse("Invalid credentials", http.StatusUnauthorized)
	}

	return s.createTokenGroup(ctx, user.Id)
}

func (s *DefaultUseService) Refresh(ctx context.Context, token tokens.Token) (*models.TokenGroup, *utils.ErrorResponse) {
	tokenId, err := uuid.Parse(token.ID)
	if err != nil {
		return nil, utils.InvalidTokenErrorResponse()
	}

	userId, err := s.tokensRepository.CheckToken(ctx, tokenId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.InvalidTokenErrorResponse()
	} else if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	err = s.tokensRepository.DeleteToken(ctx, tokenId)
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	return s.createTokenGroup(ctx, userId)
}

func NewDefaultUserService(userRepository repositories.UserRepository, tokenRepository repositories.TokenRepository, authenticator *tokens.JWTAuthenticator) *DefaultUseService {
	return &DefaultUseService{
		userRepository:   userRepository,
		tokensRepository: tokenRepository,
		authenticator:    authenticator,
	}
}
