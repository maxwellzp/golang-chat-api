package config

import (
	"github.com/joho/godotenv"
	"log"
)

type DbConfig struct {
}

type ServerConfig struct {
}

type Config struct {
	Db     DbConfig
	Server ServerConfig
}

func Load() *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	return &Config{}
}
