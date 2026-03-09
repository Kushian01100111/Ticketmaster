package auth

import (
	"context"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/auth"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

type AuthService interface {
	Login(ctx context.Context) error
	Refresh(ctx context.Context) error
	Logout(ctx context.Context) error
	SignupRequest(ctx context.Context) error
	
	SignupVeriry(ctx context.Context) error
	
	LoginRequest(ctx context.Context) error

	LoginVerify(ctxt context.Context) error
}

type authService struct {
	authRepo   repository.AuthRepository
	userRepo   repository.UserRepository
	refreshTTL time.Duration
	jwt        *auth.JWTManager
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
