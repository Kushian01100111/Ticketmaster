package auth

import (
	"context"
	"errors"
	"os/user"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/app/email"
	"github.com/Kushian01100111/Tickermaster/internal/domain/auth"
	"github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

var (
	ErrInvalidCreadentials = errors.New("invalid creadentials")
	ErrMethodNotAllowed    = errors.New("this method of sign in is no available for this user")
)

type LoginParams struct {
	Email    string
	Password string
}

type VerifyParams struct {
	Email string
	Code  string
}

type AuthService interface {
	Login(ctx context.Context, params LoginParams) (*dto.AuthResult, error)
	Refresh(ctx context.Context, refresh string) (*dto.AuthResult, error)
	Logout(ctx context.Context, refresh string) error
	SignupRequest(ctx context.Context, email string) error
	SignupVeriry(ctx context.Context, param VerifyParams) (*dto.AuthResult, error)
	LoginRequest(ctx context.Context, email string) error
	LoginVerify(ctx context.Context, param VerifyParams) (*dto.AuthResult, error)
}

type authService struct {
	otpRepo    repository.OTPRepo
	authRepo   repository.AuthRepository
	userRepo   repository.UserRepository
	mailer     email.EmailSender
	jwt        *auth.JWTManager
	otpTTL     time.Duration
	refreshTTL time.Duration
}

type AuthConfig struct {
	OTPTTL     time.Duration
	RefreshTTL time.Duration
}

func NewAuthService(
	otpRepo repository.OTPRepo,
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	mailer email.EmailSender,
	jwt *auth.JWTManager,
	config AuthConfig) AuthService {

	otpTTL := config.OTPTTL
	if otpTTL <= 0 {
		otpTTL = 10 * time.Minute
	}

	refreshTTL := config.RefreshTTL
	if refreshTTL <= 0 {
		refreshTTL = 30 * 24 * time.Hour
	}

	return &authService{
		authRepo:   authRepo,
		userRepo:   userRepo,
		mailer:     mailer,
		jwt:        jwt,
		otpTTL:     otpTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Login(ctx context.Context, params LoginParams) (*dto.AuthResult, error) {
	email := normalizeEmail(params.Email)
	pass := strings.TrimSpace(params.Password)

	if email == "" || pass == "" {
		return nil, ErrInvalidCreadentials
	}

	user, err := s.userRepo.GetByEmail(email, ctx)
	if err != nil || user == nil {
		return nil, err
	}

	if err := challengePass(user, pass); err != nil {
		return nil, err
	}

}

func (s *authService) Refresh(ctx context.Context, refresh string) (*dto.AuthResult, error) {

}

func (s *authService) Logout(ctx context.Context, refresh string) error {

}

func (s *authService) SignupRequest(ctx context.Context, email string) error {

}

func (s *authService) SignupVeriry(ctx context.Context, param VerifyParams) (*dto.AuthResult, error) {

}

func (s *authService) LoginRequest(ctx context.Context, email string) error {

}

func (s *authService) LoginVerify(cxt context.Context, param VerifyParams) (*dto.AuthResult, error) {

}

func (s *authService) challengePass(user user.User, pass string) error {
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return ErrMethodNotAllowed
	}

	if err := s.
}

///

func normalizeEmail(email string) string {
	s := strings.ToLower(strings.TrimSpace())
	if s == "" || !strings.Contains(email, "@") {
		return ""
	}
	return s
}


