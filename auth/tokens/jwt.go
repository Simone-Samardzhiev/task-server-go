package tokens

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"server/config"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

// TokenType is a custom type for token that confirm to int.
type TokenType int

const (
	// RefreshTokenType is a type of token used to revalidate tokens.
	RefreshTokenType TokenType = iota
	// AccessTokenType is a type of token used to gain access to data.
	AccessTokenType
)

// ClaimsKey is a custom type for setting the claims in context.
// It conforms to string.
type ClaimsKey string

// JWTClaimsKey is a key used to get claims.
const JWTClaimsKey = "claims"

// Token struct holds token data.
type Token struct {
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTAuthenticator used to authenticate user with JWT.
type JWTAuthenticator struct {
	secret []byte
	issuer string
}

// CreateRefreshToken will create a new [Token] with set type of [RefreshTokenType]
func (a *JWTAuthenticator) CreateRefreshToken(tokenID uuid.UUID, exp time.Time) (string, error) {
	token := Token{
		TokenType: RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.issuer,
		},
	}

	hashToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString(a.secret)
	return hashToken, err
}

// CreateAccessToken will create a new [Token] with set type of [AccessTokenType]
func (a *JWTAuthenticator) CreateAccessToken(userId int, exp time.Time) (string, error) {
	token := Token{
		TokenType: AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userId),
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.issuer,
		},
	}

	hashToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString(a.secret)
	return hashToken, err
}

// VerifyToken will verify if the jwt token in [Token] type and check its type.
func (a *JWTAuthenticator) VerifyToken(tokenString string, tokenType TokenType) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Token)
	if !ok || !token.Valid || claims.TokenType != tokenType {
		return nil, errors.New("invalid token")
	}

	return claims, err
}

// Middleware method return a middleware for jwt token.
// The middleware will extract the token and pass it by the context.
// If the token is not valid the middleware will drop the request.
func (a *JWTAuthenticator) Middleware(next http.Handler, tokenType TokenType) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.HandleErrorResponse(w, utils.NewErrorResponse("Missing bearer token", http.StatusUnauthorized))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := a.VerifyToken(token, tokenType)
		if err != nil {
			utils.HandleErrorResponse(w, utils.NewErrorResponse("Invalid token", http.StatusUnauthorized))
			return
		}

		if claims.TokenType != tokenType {
			utils.HandleErrorResponse(w, utils.NewErrorResponse("Invalid token type", http.StatusUnauthorized))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTClaimsKey, *claims)))
	})
}

func NewJWTAuthenticator(conf *config.AuthConfig) *JWTAuthenticator {
	return &JWTAuthenticator{conf.JwtSecret, conf.JwtIssuer}
}
