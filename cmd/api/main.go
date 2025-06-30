package main

import (
	"github.com/maxwellzp/golang-chat-api/internal/auth"
	"github.com/maxwellzp/golang-chat-api/internal/config"
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

	// REST API Handlers
	router := http.NewServeMux()
	auth.NewAuthHandler(router)
	room.NewRoomHandler(router)
	message.NewMessageHandler(router)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	sugar.Infof("Server running on port :%s", cfg.Server.Port)
	server.ListenAndServe()
}
