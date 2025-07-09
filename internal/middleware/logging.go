package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/maxwellzp/golang-chat-api/internal/contextkey"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			userID := "anonymous"
			if val := r.Context().Value(contextkey.UserID); val != nil {
				if id, ok := val.(int64); ok {
					userID = formatUserID(id)
				}
			}

			log.Infow("HTTP Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"user_id", userID,
				"user_agent", r.UserAgent(),
				"ip", r.RemoteAddr,
				"duration", time.Since(start).String(),
			)
		})
	}
}

func formatUserID(id int64) string {
	return "user#" + strconv.FormatInt(id, 10)
}
