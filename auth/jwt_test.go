package auth

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncode(t *testing.T) {
	id := uuid.New()
	sub := uuid.New()

	// Check if GenerateTokenGroup will return an error.
	group, err := DefaultJWTService.GenerateTokenGroup(&id, &sub)
	if err != nil {
		t.Error(err)
	}

	t.Log(group)
}

func TestDecode(t *testing.T) {
	id := uuid.New()
	sub := uuid.New()

	group, err := DefaultJWTService.GenerateTokenGroup(&id, &sub)
	if err != nil {
		t.Error(err)
	}
	t.Log(group)

	// Check if parse token will return an error if the token is of type access.
	accessToken, err := DefaultJWTService.ParseToken(group.AccessToken)
	if err != nil {
		t.Error(err)
	}

	// Check if the access token id match with parse one.
	if accessToken.Subject != sub.String() {
		t.Error("access token subject does not match")
	}

	// Check if parse token will return an error if the token is of type refresh.
	refreshToken, err := DefaultJWTService.ParseToken(group.RefreshToken)
	if err != nil {
		t.Error(err)
	}

	// Check if the refresh token id and subject match the parse one.
	if refreshToken.Subject != sub.String() || refreshToken.ID != id.String() {
		t.Error("refresh token id or sub does not match")
	}
}

func TestJWTMiddlewareForRefreshToken(t *testing.T) {
	id := uuid.New()
	sub := uuid.New()
	tokenGroup, err := DefaultJWTService.GenerateTokenGroup(&id, &sub)

	if err != nil {
		t.Error(err)
	}

	// Create a handler that expect a refresh token.
	handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the claims and check them.
		token, ok := r.Context().Value(TokenKey).(*CustomClaims)
		if !ok {
			t.Error("failed to get custom claims")
			return
		}

		// Check if the id, subject and type match.
		if token.ID != id.String() || token.Subject != sub.String() || token.Type != RefreshToken {
			t.Error("token id or sub does not match or the type is not refresh token")
		}

		w.WriteHeader(http.StatusOK)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer "+tokenGroup.RefreshToken)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestJWTMiddlewareForAccessToken(t *testing.T) {
	id := uuid.New()
	sub := uuid.New()
	tokenGroup, err := DefaultJWTService.GenerateTokenGroup(&id, &sub)
	if err != nil {
		t.Error(err)
	}

	// Create a handler that expect an access token.
	handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value(TokenKey).(*CustomClaims)
		if !ok {
			t.Error("failed to get custom claims")
		}

		// Check if the token subject and type match.
		if token.Subject != sub.String() {
			t.Error("token id does not match")
		}

		w.WriteHeader(http.StatusOK)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer "+tokenGroup.AccessToken)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
