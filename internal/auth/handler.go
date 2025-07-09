package auth

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/validatorx"
	"net/http"
)

type AuthHandler struct {
	authService *AuthService
	validator   *validatorx.Validator
	logger      *logger.Logger
}

func NewAuthHandler(authService *AuthService, validator *validatorx.Validator, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Warnw("Login request JSON decode failed",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			h.logger.Infow("Login request validation failed",
				"email", req.Email,
				"errors", err,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		user, token, err := h.authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			h.logger.Warnw("Login failed",
				"email", req.Email,
				"error", err,
			)
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		h.logger.Infow("User logged in successfully",
			"user_id", user.ID,
			"email", req.Email,
		)
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
			h.logger.Warnw("Register request JSON decode failed",
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := h.validator.Validate(&req); err != nil {
			h.logger.Infow("Register request validation failed",
				"email", req.Email,
				"errors", err,
			)
			httpx.WriteValidationError(w, err)
			return
		}

		user, err := h.authService.Register(r.Context(), req.Username, req.Email, req.Password)
		if err != nil {
			h.logger.Warnw("User registration failed",
				"email", req.Email,
				"error", err,
			)
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.logger.Infow("User registered successfully",
			"user_id", user.ID,
			"email", req.Email,
		)
		httpx.WriteJSON(w, http.StatusCreated, user)
	}
}
