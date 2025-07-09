package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/maxwellzp/golang-chat-api/internal/contextkey"
	"github.com/maxwellzp/golang-chat-api/internal/httpx"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"net/http"
	"strings"
)

func JWT(secret string, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warnw("Missing Authorization header",
					"path", r.URL.Path,
				)
				httpx.WriteError(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Warnw("Invalid Authorization header format",
					"header", authHeader,
				)
				httpx.WriteError(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}

			tokenStr := parts[1]
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Errorw("Unexpected signing method",
						"method", token.Method.Alg(),
					)
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				log.Warnw("Invalid or expired token",
					"error", err,
				)
				httpx.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Extract user ID from claims
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				log.Errorw("user_id missing or invalid in token",
					"claims", claims,
				)
				httpx.WriteError(w, http.StatusUnauthorized, "Invalid user_id in token")
				return
			}
			log.Infow("Authenticated request",
				"user_id", int64(userIDFloat),
				"path", r.URL.Path,
			)
			ctx := context.WithValue(r.Context(), contextkey.UserID, int64(userIDFloat))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
