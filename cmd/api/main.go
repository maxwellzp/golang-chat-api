package main

import (
	"context"
	"github.com/maxwellzp/golang-chat-api/internal/auth"
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/message"
	"github.com/maxwellzp/golang-chat-api/internal/room"
	"log"
	"net/http"
)

func main() {
	sugar, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	cfg := config.Load(sugar)
	sugar.Infof("%+v", cfg)

	dbInstance, err := db.NewDb(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}

	// Database Repositories
	authRepo := auth.NewAuthRepository(dbInstance)
	roomRepo := room.NewRoomRepository(dbInstance)
	messageRepo := message.NewMessageRepository(dbInstance)

	// Business logic services
	authService := auth.NewAuthService(authRepo)
	roomService := room.NewRoomService(roomRepo)
	messageService := message.NewMessageService(messageRepo)

	// REST API Handlers
	router := http.NewServeMux()
	auth.NewAuthHandler(router, authService)
	room.NewRoomHandler(router, roomService)
	message.NewMessageHandler(router, messageService)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	sugar.Infof("Server running on port :%s", cfg.Server.Port)
	server.ListenAndServe()
}
