package handlers

import (
	"encoding/json"
	"net/http"
	"server/models"
	"server/services"
	"server/utils"
)

// UserHandler interface handles user requests.
type UserHandler interface {
	Register() http.HandlerFunc

	Login() http.HandlerFunc
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

	}
}

func NewDefaultUserHandler(userRepository services.UserService) *DefaultUserHandler {
	return &DefaultUserHandler{
		userService: userRepository,
	}
}
