package auth

import (
	"context"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/app/email"
	"github.com/Kushian01100111/Tickermaster/internal/domain/auth"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

type LoginParams struct {
	Email    string
	Password string
}

type AuthService interface {
	Login(ctx context.Context) (dto.AuthResult, error)
	Refresh(ctx context.Context) (dto.AuthResult, error)
	Logout(ctx context.Context) error
	SignupRequest(ctx context.Context) error
	SignupVeriry(ctx context.Context) (dto.AuthResult, error)
	LoginRequest(ctx context.Context) error
	LoginVerify(ctxt context.Context) (dto.AuthResult, error)
}

type authService struct {
	OTPRepo    repository.OTPRepo
	authRepo   repository.AuthRepository
	userRepo   repository.UserRepository
	mailer     email.EmailSender
	jwt        *auth.JWTManager
	otpTTL     time.Duration
	refreshTTL time.Duration
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, mailer email.EmailSender, jwt *auth.JWTManager, otpTTL time.Duration, refreshTTL time.Duration) AuthService {
	return &authService{
		authRepo:   authRepo,
		userRepo:   userRepo,
		refreshTTL: refreshTTL,
		mailer:     mailer,
		jwt:        jwt,
		otpTTL:     otpTTL,
		refreshTTL: refreshTTL,
	}
}
