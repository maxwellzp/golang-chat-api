package auth

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"github.com/maxwellzp/golang-chat-api/internal/validatorx"
	"net/http"
)

type AuthHandler struct {
	authService *AuthService
	validator   *validatorx.Validator
}

func NewAuthHandler(authService *AuthService, validator *validatorx.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			httpx.WriteValidationError(w, err)
			return
		}

		_, token, err := h.authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		resp := LoginResponse{
			Token: token,
		}
		httpx.WriteJSON(w, http.StatusOK, resp)
	}
}
func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			httpx.WriteValidationError(w, err)
			return
		}

		u, err := h.authService.Register(r.Context(), req.Username, req.Email, req.Password)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, u)
	}
}
