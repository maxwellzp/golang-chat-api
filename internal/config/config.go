package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type ApplicationConfig struct {
	AppEnv string
}

type DbConfig struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Application ApplicationConfig
	Db          DbConfig
	Server      ServerConfig
}

func Load(logger *zap.SugaredLogger) *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// Environment variables can also be set via terminal, bash script ...
		// APP_ENV=dev go run cmd/api/main.go
		logger.Warnw("No .env file found")
	}

	return &Config{
		Application: ApplicationConfig{
			AppEnv: getEnv(logger, "APP_ENV", "prod"),
		},
		Db: DbConfig{
			PostgresUser:     mustGetEnv(logger, "POSTGRES_USER"),
			PostgresPassword: mustGetEnv(logger, "POSTGRES_PASSWORD"),
			PostgresDatabase: mustGetEnv(logger, "POSTGRES_DB"),
			PostgresHost:     mustGetEnv(logger, "POSTGRES_HOST"),
			PostgresPort:     getEnv(logger, "POSTGRES_PORT", "5432"),
		},
		Server: ServerConfig{
			Port: getEnv(logger, "SERVER_PORT", "8080"),
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
