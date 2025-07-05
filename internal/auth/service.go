package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/maxwellzp/golang-chat-api/internal/user"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type AuthService struct {
	userRepository *user.UserRepository
	jwtSecret      string
}

func NewAuthService(userRepository *user.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepository: userRepository, jwtSecret: jwtSecret}
}

func (as *AuthService) Register(ctx context.Context, username, email, password string) (*user.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)

	if username == "" || email == "" || password == "" {
		return nil, errors.New("missing required fields")
	}
	existingUser, err := as.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := as.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil

}

func (as *AuthService) Login(ctx context.Context, email string, password string) (*user.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	existingUser, err := as.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existingUser == nil {
		return nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	existingUser.Password = ""
	return existingUser, tokenString, nil
}
