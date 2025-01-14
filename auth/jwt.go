package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"server/config"
	"server/utils"
	"strings"
	"time"
)

// TokenType is used to set the type of token.
type TokenType int

const (
	// RefreshToken is type of token used to refresh access.
	RefreshToken TokenType = 1
	// AccessToken is type of token used to access information.
	AccessToken TokenType = 2
)

// CustomClaims type struct is used to create the jwt claims.
type CustomClaims struct {
	Type TokenType `json:"type"`
	jwt.RegisteredClaims
}

// TokenGroup type claims used to hold the access token and refresh token.
type TokenGroup struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWTService type struct used to create a service that will manage jwt.
type JWTService struct {
	secret []byte
}

// GenerateTokenGroup will return [TokenGroup] setting the tokens id and subject.
func (s *JWTService) GenerateTokenGroup(id, sub *uuid.UUID, exp *time.Time) (*TokenGroup, error) {
	// Creating the refresh token.
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Type: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Subject:   sub.String(),
			ExpiresAt: jwt.NewNumericDate(*exp),
		},
	})

	// Hashing the refresh token.
	refreshTokenString, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	//Creating the access token.
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Type: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})

	// Hashing the access token.
	accessTokenString, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &TokenGroup{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil

}

// ParseToken will parse a string into [CustomClaims].
func (s *JWTService) ParseToken(tokenString string) (*CustomClaims, error) {
	// Parse the token.
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the hash method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})

	if err != nil {
		return nil, utils.UnauthorizedErr
	}

	// Getting the claims.
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, utils.UnauthorizedErr
}

// DefaultJWTService is [JWTService] used by the app.
var DefaultJWTService = newJWTService(config.Envs.JWTSecret)

// newJWTService will create a new service setting the secret.
func newJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

// Key is the type used to create context key.
type Key string

// TokenKey is the key that is used to retrieve token from context.
const TokenKey Key = "token"

// JWTMiddleware used to wrap handlers if they need to be access by token.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting the header.
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		// Parsing the token.
		token, err := DefaultJWTService.ParseToken(tokenString)

		// Any error means the token is invalid.
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		// Putting the token in a context so the next handler can access it.
		ctx := context.WithValue(r.Context(), TokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
