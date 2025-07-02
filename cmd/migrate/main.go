package main

import (
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"go.uber.org/zap"
	"log"
)

func main() {
	// Temporary logger to catch early config errors
	tempLogger := zap.NewExample().Sugar()

	// Load config
	cfg := config.Load(tempLogger)

	// Run migrations from /migrations
	if err := db.RunMigrations(cfg); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	tempLogger.Infow("Migrations applied successfully")
}
