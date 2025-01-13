package auth

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncode(t *testing.T) {
	refreshClaims := NewRefreshTokenClaims(uuid.New(), uuid.New())

	// Check if Encode function works for RefreshTokenClaims.
	_, err := Encode(refreshClaims)
	if err != nil {
		t.Error(err)
	}

	// Check if Encode function works for AccessTokenClaims.
	accessClaims := NewAccessTokenClaims(uuid.New())
	_, err = Encode(accessClaims)
	if err != nil {
		t.Error(err)
	}
}

func TestRefreshTokenMiddleware(t *testing.T) {
	refreshClaims := NewRefreshTokenClaims(uuid.New(), uuid.New())
	token, err := Encode(refreshClaims)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/refresh_token", nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := RefreshTokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(RefreshTokenClaimsKey).(RefreshTokenClaims)

		// Check if the claims match.
		if claims.ID != refreshClaims.ID || claims.Subject != refreshClaims.Subject {
			t.Errorf("RefreshTokenMiddleware got claims %+v, expected %+v", claims, refreshClaims)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	// Check if the status code is OK.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAccessTokenMiddleware(t *testing.T) {
	accessClaims := NewAccessTokenClaims(uuid.New())
	token, err := Encode(accessClaims)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/access_token", nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := AccessTokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(AccessTokenClaimsKey).(AccessTokenClaims)

		// Check if the claims match.
		if claims.ID != accessClaims.ID || claims.Subject != accessClaims.Subject {
			t.Errorf("AccessTokenMiddleware got claims %+v, expected %+v", claims, accessClaims)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Check if the status code is OK.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
