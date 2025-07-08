package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/maxwellzp/golang-chat-api/internal/auth"
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/message"
	appMiddleware "github.com/maxwellzp/golang-chat-api/internal/middleware"
	"github.com/maxwellzp/golang-chat-api/internal/room"
	"github.com/maxwellzp/golang-chat-api/internal/user"
	validatorx "github.com/maxwellzp/golang-chat-api/internal/validatorx"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
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
		tempLogger.Fatalw("Failed to initialize logger",
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

	// Validator
	val := validatorx.NewValidator()

	// REST API Handlers
	authHandler := auth.NewAuthHandler(authService, val)
	roomHandler := room.NewRoomHandler(roomService, val)
	messageHandler := message.NewMessageHandler(messageService, val)

	// Middleware
	jwtMiddleWare := appMiddleware.JWT(cfg.Auth.JwtSecret)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check (public)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Auth (public)
	r.Post("/login", authHandler.Login())
	r.Post("/register", authHandler.Register())

	// Messages (all protected)
	r.Route("/messages", func(r chi.Router) {
		r.Use(jwtMiddleWare)

		r.Post("/create", messageHandler.Create())
		r.Patch("/update/{id}", messageHandler.Update())
		r.Delete("/delete/{id}", messageHandler.Delete())
		r.Get("/{id}", messageHandler.GetByID())
		r.Get("/list", messageHandler.List())
	})

	// Rooms
	r.Route("/rooms", func(r chi.Router) {
		// Public endpoints
		r.Get("/list", roomHandler.List())
		r.Get("/{id}", roomHandler.GetByID())

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(jwtMiddleWare)
			r.Post("/create", roomHandler.Create())
			r.Patch("/update/{id}", roomHandler.Update())
			r.Delete("/delete/{id}", roomHandler.Delete())
		})
	})

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	log.Infow("Server running",
		"port", cfg.Server.Port,
		"env", cfg.Application.AppEnv,
	)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw("Failed to start server",
				"port", cfg.Server.Port,
				"env", cfg.Application.AppEnv,
				"err", err,
			)
		}
	}()
	<-shutdownCtx.Done()
	log.Infow("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("Failed to shutdown server",
			"err", err,
		)
	} else {
		log.Infow("Server gracefully stopped")
	}
}
