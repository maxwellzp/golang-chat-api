package auth

import (
	"context"
	"github.com/maxwellzp/golang-chat-api/internal/db"
	"github.com/maxwellzp/golang-chat-api/internal/user"
)

type AuthRepository struct {
	Database *db.Db
}

func NewAuthRepository(database *db.Db) *AuthRepository {
	return &AuthRepository{Database: database}
}

func (r *AuthRepository) Create(ctx context.Context, user *user.User) error {
	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}
