package otp

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/otp"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

var (
	ErrPurpuse           = errors.New("invalid purpuse type")
	ErrEmailRequired     = errors.New("email is required")
	ErrCodeHasRequired   = errors.New("codeHas is required")
	ErrExpiredAtRequired = errors.New("expiredAt is required")
)

type OTPParams struct {
	Email     string
	Purpuse   string
	CodeHash  string
	ExpiresAt time.Time
	Attempts  int32
	CreatedAt time.Time
}

type OTPService interface {
	CreateOrReplace(ctx context.Context, ch OTPParams) (string, error)
	GetActiveByEmail(ctx context.Context, email string) (*otp.OTPChallange, error)
	IncAttempts(ctx context.Context, email string) error
	Consume(ctx context.Context, email string) error
}

type otpService struct {
	otpRepo  repository.OTPRepo
	userRepo repository.UserRepository
}

func NewOTPService(otp repository.OTPRepo, user repository.UserRepository) OTPService {
	return &otpService{
		otpRepo:  otp,
		userRepo: user,
	}
}

func (o *otpService) CreateOrReplace(ctx context.Context, ch OTPParams) (string, error) {
	if err := validateParams(ch); err != nil {
		return "", err
	}

	now := time.Now()
	if ch.CreatedAt.IsZero() {
		ch.CreatedAt = now
	}

	str, err := o.otpRepo.CreateOrReplace(ctx, otp.OTPChallange{
		Email:     ch.Email,
		Purpuse:   ch.Purpuse,
		CodeHash:  ch.CodeHash,
		ExpiresAt: ch.ExpiresAt,
		Attempts:  ch.Attempts,
		CreatedAt: now,
	})
	if err != nil {
		return "", err
	}

	return str, nil
}
func (o *otpService) GetActiveByEmail(ctx context.Context, email string) (*otp.OTPChallange, error)
func (o *otpService) IncAttempts(ctx context.Context, email string) error
func (o *otpService) Consume(ctx context.Context, email string) error

//

func validateParams(otp OTPParams) error {
	if strings.TrimSpace(otp.Email) == "" {
		return ErrEmailRequired
	}

	if err := validatePurpuse(otp.Purpuse); err != nil {
		return ErrPurpuse
	}

	if strings.TrimSpace(otp.CodeHash) == "" {
		return ErrCodeHasRequired
	}

	if otp.ExpiresAt.IsZero() {
		return ErrExpiredAtRequired
	}

	return nil
}

type Purpuse string

const (
	PurpuseLogin  Purpuse = "login"
	PurpuseSignUp Purpuse = "signUp"
)

func validatePurpuse(str string) error {
	switch Purpuse(str) {
	case PurpuseLogin, PurpuseSignUp:
		return nil
	default:
		return ErrPurpuse
	}
}
