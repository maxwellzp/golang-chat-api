package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type ApplicationConfig struct {
	AppEnv string
}

type AuthConfig struct {
	JwtSecret string
}

type DbConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func (c *DbConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Application ApplicationConfig
	Db          DbConfig
	Server      ServerConfig
	Auth        AuthConfig
}

func Load(logger *zap.SugaredLogger) *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// Environment variables can also be set via terminal, bash script ...
		// APP_ENV=dev go run cmd/server/main.go
		logger.Warnw("No .env file found")
	}

	return &Config{
		Application: ApplicationConfig{
			AppEnv: getEnv(logger, "APP_ENV", "prod"),
		},
		Db: DbConfig{
			User:     mustGetEnv(logger, "POSTGRES_USER"),
			Password: mustGetEnv(logger, "POSTGRES_PASSWORD"),
			Database: mustGetEnv(logger, "POSTGRES_DB"),
			Host:     mustGetEnv(logger, "POSTGRES_HOST"),
			Port:     getEnv(logger, "POSTGRES_PORT", "5432"),
		},
		Server: ServerConfig{
			Port: getEnv(logger, "SERVER_PORT", "8080"),
		},
		Auth: AuthConfig{
			JwtSecret: mustGetEnv(logger, "JWT_SECRET"),
		},
	}
}

func getEnv(logger *zap.SugaredLogger, key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	logger.Infow("Using default value for env variable",
		"environment variable", key,
		"default", defaultVal,
	)
	return defaultVal
}

func mustGetEnv(logger *zap.SugaredLogger, key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		logger.Fatalw("Required environment variable is not set",
			"key", key,
		)
	}
	return value
}
