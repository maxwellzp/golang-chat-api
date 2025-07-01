package auth

type AuthService struct {
	AuthRepository *AuthRepository
}

func NewAuthService(authRepository *AuthRepository) *AuthService {
	return &AuthService{AuthRepository: authRepository}
}
