package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"log"
)

type Db struct {
	*sql.DB
}

func NewDb(ctx context.Context, cfg *config.Config) (*Db, error) {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.PostgresHost,
		cfg.Db.PostgresPort,
		cfg.Db.PostgresUser,
		cfg.Db.PostgresPassword,
		cfg.Db.PostgresDatabase,
	)
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if pingErr := db.PingContext(ctx); pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	return &Db{db}, nil
}
