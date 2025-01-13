package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"server/config"
	"strings"
	"time"
)

// RefreshTokenClaims used to create refresh token.
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
}

// NewRefreshTokenClaims used to create a new token with specified id and subject.
func NewRefreshTokenClaims(id, sub *uuid.UUID) RefreshTokenClaims {
	return RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Subject:   sub.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 14)),
		},
	}
}

// newRefreshTokenClaimsFromToken used to parse the token and return the claims or an error.
func newRefreshTokenClaimsFromToken(tokenString *string) (RefreshTokenClaims, error) {
	var claims RefreshTokenClaims
	token, err := jwt.ParseWithClaims(*tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return RefreshTokenClaims{}, err
	}

	if !token.Valid {
		return RefreshTokenClaims{}, errors.New("invalid token")
	}

	return claims, nil
}

// AccessTokenClaims used to create access token.
type AccessTokenClaims struct {
	jwt.RegisteredClaims
}

// NewAccessTokenClaims used to create a new token with subject.
func NewAccessTokenClaims(sub *uuid.UUID) AccessTokenClaims {
	return AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}
}

// newAccessTokenClaimsFromToken used to parse the token and return the claims or an error.
func newAccessTokenClaimsFromToken(tokenString *string) (AccessTokenClaims, error) {
	var claims AccessTokenClaims
	token, err := jwt.ParseWithClaims(*tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return AccessTokenClaims{}, err
	}
	if !token.Valid {
		return AccessTokenClaims{}, errors.New("invalid token")
	}

	return claims, nil
}

// ContextKey is a type used to declare keys for middleware.
type ContextKey string

// RefreshTokenClaimsKey is a key used to get refresh claims.
// They will only be sent if the handler is wrapped in the middleware.
// If the token is not valid the middleware will automatically return.
const RefreshTokenClaimsKey ContextKey = "rtc"

// AccessTokenClaimsKey is a key used to get refresh claims.
// They will only be sent if the handler is wrapped in the middleware.
// If the token is not valid the middleware will automatically return.
const AccessTokenClaimsKey ContextKey = "atc"

// RefreshTokenMiddleware will wrap a handler in middleware that will check for valid refresh token.
func RefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims, err := newRefreshTokenClaimsFromToken(&header)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RefreshTokenClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AccessTokenMiddleware will wrap a handler in middleware that will check for valid access token.
func AccessTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims, err := newAccessTokenClaimsFromToken(&header)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), AccessTokenClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Encode will sign tokens.
func Encode(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(config.Envs.JWTSecret))
}
