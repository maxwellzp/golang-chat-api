package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/maxwellzp/golang-chat-api/internal/config"
)

type Db struct {
	*sql.DB
}

func NewDb(ctx context.Context, cfg *config.Config) (*Db, error) {
	db, err := sql.Open("postgres", cfg.Db.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	if pingErr := db.PingContext(ctx); pingErr != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return &Db{DB: db}, nil
}
