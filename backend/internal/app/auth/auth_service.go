package auth

import "github.com/Kushian01100111/Tickermaster/internal/repository"

type AuthService interface{}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
