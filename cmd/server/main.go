package main

import (
	"context"
	"github.com/maxwellzp/golang-chat-api/internal/auth"
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/message"
	"github.com/maxwellzp/golang-chat-api/internal/room"
	"github.com/maxwellzp/golang-chat-api/internal/user"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	// Temporary logger to catch early config errors
	tempLogger := zap.NewExample().Sugar()

	// Load config
	cfg := config.Load(tempLogger)

	// Instantiate proper Zap logger based on APP_ENV
	log, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalw("Failed to initialize logger",
			"err", err,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	dbInstance, err := db.NewDb(ctx, cfg)
	if err != nil {
		log.Fatalw("Failed to initialize db",
			"err", err,
		)
	}
	log.Debugw("Database connection established")

	// Instantiate database repositories
	userRepo := user.NewUserRepository(dbInstance)
	roomRepo := room.NewRoomRepository(dbInstance)
	messageRepo := message.NewMessageRepository(dbInstance)

	// Instantiate business logic services
	authService := auth.NewAuthService(userRepo, cfg.Auth.JwtSecret)
	roomService := room.NewRoomService(roomRepo)
	messageService := message.NewMessageService(messageRepo)

	// Bind REST API Handlers to ServeMux
	router := http.NewServeMux()
	auth.NewAuthHandler(router, authService)
	room.NewRoomHandler(router, roomService)
	message.NewMessageHandler(router, messageService)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	log.Infow("Server running",
		"port", cfg.Server.Port,
		"env", cfg.Application.AppEnv,
	)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalw("Failed to start server",
			"port", cfg.Server.Port,
			"env", cfg.Application.AppEnv,
			"err", err,
		)
	}
}
