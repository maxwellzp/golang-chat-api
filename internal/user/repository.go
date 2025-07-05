package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"time"
)

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{database: database}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password, created_at) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at`

	row := r.database.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, time.Now())
	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	row := r.database.QueryRowContext(ctx, query, email)

	user := &User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
