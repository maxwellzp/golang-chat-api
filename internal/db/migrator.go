package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/maxwellzp/golang-chat-api/internal/config"
)

func RunMigrations(cfg *config.Config) error {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.PostgresHost,
		cfg.Db.PostgresPort,
		cfg.Db.PostgresUser,
		cfg.Db.PostgresPassword,
		cfg.Db.PostgresDatabase,
	)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return errors.New("failed to connect to DB: " + err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.New("failed to create DB driver: " + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return errors.New("failed to initialize migrate: " + err.Error())
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.New("migration error: " + err.Error())
	}
	return nil
}
