package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/maxwellzp/golang-chat-api/internal/logger"
	"github.com/maxwellzp/golang-chat-api/internal/user"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type AuthService struct {
	userRepository *user.UserRepository
	jwtSecret      string
	logger         *logger.Logger
}

func NewAuthService(userRepository *user.UserRepository, jwtSecret string, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		logger:         logger,
	}
}

func (as *AuthService) Register(ctx context.Context, username, email, password string) (*user.User, error) {
	// Normalization
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)

	existingUser, err := as.userRepository.FindByEmail(ctx, email)
	if err != nil {
		as.logger.Errorw("Failed to check existing user",
			"error", err,
		)
		return nil, err
	}
	if existingUser != nil {
		as.logger.Warnw("Registration failed: email already in use",
			"email", email,
		)
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		as.logger.Errorw("Password hashing failed",
			"error", err,
		)
		return nil, err
	}

	u := &user.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := as.userRepository.Create(ctx, u); err != nil {
		as.logger.Errorw("User creation failed",
			"error", err,
		)
		return nil, err
	}
	u.Password = ""
	return u, nil
}

func (as *AuthService) Login(ctx context.Context, email string, password string) (*user.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	existingUser, err := as.userRepository.FindByEmail(ctx, email)
	if err != nil {
		as.logger.Errorw("Login failed: DB error",
			"email", email,
			"error", err,
		)
		return nil, "", err
	}
	if existingUser == nil {
		as.logger.Warnw("Login failed: user not found",
			"email", email,
		)
		return nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		as.logger.Warnw("Login failed: incorrect password",
			"email", email,
		)
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		as.logger.Errorw("JWT generation failed",
			"email", email,
			"error", err,
		)
		return nil, "", err
	}

	existingUser.Password = ""
	return existingUser, tokenString, nil
}
