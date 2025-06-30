package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

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
	Db     DbConfig
	Server ServerConfig
}

func Load(logger *zap.SugaredLogger) *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warnw("No .env file found")
	}
	return &Config{
		Db: DbConfig{
			PostgresUser:     os.Getenv("POSTGRES_USER"),
			PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
			PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
			PostgresHost:     os.Getenv("POSTGRES_HOST"),
			PostgresPort:     os.Getenv("POSTGRES_PORT"),
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
	}
}
