package main

import (
	"fmt"
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ApiHandler struct{}

func NewApiHandler(router *http.ServeMux) {
	handler := &ApiHandler{}
	router.HandleFunc("/api", handler.handleRequest())
	router.HandleFunc("/api2", handler.handleRequest2())
}

func (h *ApiHandler) handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received1")
		fmt.Fprintf(w, "Response from API1")
	}
}

func (h *ApiHandler) handleRequest2() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received2")
		fmt.Fprintf(w, "Response from API2")
	}
}

func main() {
	url := "localhost:8080"
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	cfg := config.Load()
	fmt.Printf("%+v\n", cfg)

	router := http.NewServeMux()
	NewApiHandler(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server running on port 8080")
	server.ListenAndServe()
}
