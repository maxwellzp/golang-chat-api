package auth

import (
	"context"
	"github.com/maxwellzp/golang-chat-api/internal/user"
)

type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) Create(ctx context.Context, user *user.User) error {
	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}
