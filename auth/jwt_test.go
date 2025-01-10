package auth

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncode(t *testing.T) {
	refreshClaims := NewRefreshTokenClaims(uuid.New(), uuid.New())
	_, err := Encode(refreshClaims)
	if err != nil {
		t.Error(err)
	}

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
		if claims.ID != refreshClaims.ID || claims.Subject != refreshClaims.Subject {
			t.Errorf("RefreshTokenMiddleware got claims %+v, expected %+v", claims, refreshClaims)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
