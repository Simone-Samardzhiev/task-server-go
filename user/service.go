package user

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"server/auth"
	"server/utils"
	"time"
)

// Service type interface manage users business logic.
type Service interface {
	// Register will register the user.
	Register(user *JsonUser) error
	// Login will create an access token and return it.
	Login(user *JsonUser) (*auth.TokenGroup, error)
	// RefreshToken will return new refresh and access token.
	RefreshToken(token *auth.CustomClaims) (*auth.TokenGroup, error)
}

// DefaultService is the default implementation of Service.
type DefaultService struct {
	repository Repository
}

func NewDefaultService(repository Repository) *DefaultService {
	return &DefaultService{repository: repository}
}

func (d *DefaultService) Register(user *JsonUser) error {
	// Fetch data if the user email is in use.
	inUse, err := d.repository.CheckUserEmail(user.Email)
	if err != nil {
		return utils.InternalServerErr
	}

	// Check if the email is already in use.
	if inUse {
		return utils.ConflictErr
	}

	// Hash the user passwords for security.
	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		return utils.InternalServerErr
	}

	// Add the user to the database.
	err = d.repository.AddUser(
		&User{
			Id:       uuid.New(),
			Email:    user.Email,
			Password: hash,
		},
	)

	if err != nil {
		return utils.InternalServerErr
	}

	return nil
}

func (d *DefaultService) Login(user *JsonUser) (*auth.TokenGroup, error) {
	// Fetching the user by email.
	fetchedUser, err := d.repository.GetUserByEmail(user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		// If the error is sql.ErrNoRows it meas user with that email doesn't exist.
		return nil, utils.UnauthorizedErr
	} else if err != nil {
		return nil, utils.InternalServerErr
	}

	// Check if the passwords match.
	if !auth.CheckPassword(user.Password, fetchedUser.Password) {
		return nil, utils.UnauthorizedErr
	}

	// Creating refresh token for the user.
	tokenId := uuid.New()
	exp := time.Now().Add(time.Hour * 24 * 30)

	group, err := auth.DefaultJWTService.GenerateTokenGroup(&tokenId, &fetchedUser.Id, &exp)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	// Delete any token the user might have, because only one should be valid.
	err = d.repository.DeleteTokenByUserId(&fetchedUser.Id)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	// Add the new token.
	err = d.repository.AddToken(&tokenId, &fetchedUser.Id, &exp)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	return group, nil
}

func (d *DefaultService) RefreshToken(token *auth.CustomClaims) (*auth.TokenGroup, error) {
	// Check if the token is of refresh token type.
	if token.Type != auth.RefreshToken {
		return nil, utils.UnauthorizedErr
	}

	// Parse the token id.
	tokenId, err := uuid.Parse(token.ID)
	if err != nil {
		return nil, utils.UnauthorizedErr
	}

	// Deleting the token.
	deleted, err := d.repository.DeleteTokenById(&tokenId)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	if !deleted {
		// If delete is false it means token with that id doesn't exist.
		return nil, utils.UnauthorizedErr
	}

	// Parse the token id.
	id, err := uuid.Parse(token.ID)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	sub, err := uuid.Parse(token.Subject)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	// Creating the group.
	exp := time.Now().Add(time.Hour * 24 * 30)
	group, err := auth.DefaultJWTService.GenerateTokenGroup(&id, &sub, &exp)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	return group, nil
}
