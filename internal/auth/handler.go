package auth

import (
	"encoding/json"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"net/http"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(router *http.ServeMux, authService *AuthService) {
	handler := &AuthHandler{
		authService: authService,
	}
	router.HandleFunc("POST /login", handler.Login())
	router.HandleFunc("POST /register", handler.Register())
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "Invalid request body")
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
		u, err := h.authService.Register(r.Context(), req.Username, req.Email, req.Password)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, u)
	}
}
