package main

import (
	"fmt"
	"github.com/maxwellzp/golang-chat-api/internal/auth"
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

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
	auth.NewAuthHandler(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server running on port 8080")
	server.ListenAndServe()
}
