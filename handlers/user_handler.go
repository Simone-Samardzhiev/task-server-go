package handlers

import (
	"encoding/json"
	"net/http"
	"server/auth/tokens"
	"server/models"
	"server/services"
	"server/utils"
)

// UserHandler interface handles user requests.
type UserHandler interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	Refresh() http.HandlerFunc
}

// DefaultUserHandler interface is the default implementation of [UserHandler]
type DefaultUserHandler struct {
	userService services.UserService
}

func (h *DefaultUserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.RegistrationsPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InvalidJson())
			return
		}

		if utils.HandlePayload(w, &payload) {
			return
		}

		response := h.userService.Register(r.Context(), payload)

		if utils.HandleErrorResponse(w, response) {
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *DefaultUserHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.LoginPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InvalidJson())
			return
		}

		tokenGroup, errorResponse := h.userService.Login(r.Context(), payload)
		if errorResponse != nil {
			utils.HandleErrorResponse(w, errorResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(tokenGroup)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InternalServerError())
			return
		}
	}
}

func (h *DefaultUserHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value(tokens.JWTClaimsKey).(tokens.Token)
		if !ok {
			utils.HandleErrorResponse(w, utils.InternalServerError())
			return
		}

		tokenGroup, errResponse := h.userService.Refresh(r.Context(), token)
		if errResponse != nil {
			utils.HandleErrorResponse(w, errResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(tokenGroup)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InternalServerError())
			return
		}
	}
}

func NewDefaultUserHandler(userRepository services.UserService) *DefaultUserHandler {
	return &DefaultUserHandler{
		userService: userRepository,
	}
}
