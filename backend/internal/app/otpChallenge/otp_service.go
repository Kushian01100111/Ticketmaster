package otpChallenge

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
	CreateOrReplace(ctx context.Context, ch OTPParams) error
	GetActiveByEmail(ctx context.Context, email, purpuse string) (*otp.OTP, error)
	IncAttempts(ctx context.Context, email, purpuse string) error
	Consume(ctx context.Context, email, purpuse string) error
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

func (o *otpService) CreateOrReplace(ctx context.Context, ch OTPParams) error {
	if err := validateParams(ch); err != nil {
		return err
	}

	now := time.Now()
	if ch.CreatedAt.IsZero() {
		ch.CreatedAt = now
	}

	err := o.otpRepo.CreateOrReplace(ctx, otp.OTP{
		Email:     strings.ToLower(strings.TrimSpace(ch.Email)),
		Purpuse:   ch.Purpuse,
		CodeHash:  ch.CodeHash,
		ExpiresAt: ch.ExpiresAt,
		Attempts:  ch.Attempts,
		CreatedAt: now,
	})
	if err != nil {
		return err
	}

	return nil
}
func (o *otpService) GetActiveByEmail(ctx context.Context, email, purpuse string) (*otp.OTP, error) {
	mail := strings.ToLower(strings.TrimSpace(email))
	purpuse = strings.TrimSpace(purpuse)

	if mail == "" {
		return nil, ErrEmailRequired
	}

	if err := validatePurpuse(purpuse); err != nil {
		return nil, err
	}

	return o.otpRepo.GetActiveByEmail(ctx, mail, purpuse)
}

func (o *otpService) IncAttempts(ctx context.Context, email, purpuse string) error {
	mail := strings.ToLower(strings.TrimSpace(email))

	if mail == "" {
		return ErrEmailRequired
	}

	return o.otpRepo.IncAttempts(ctx, mail, purpuse)
}
func (o *otpService) Consume(ctx context.Context, email, purpuse string) error {
	mail := strings.ToLower(strings.TrimSpace(email))
	if mail == "" {
		return ErrEmailRequired
	}

	now := time.Now()
	return o.otpRepo.Consume(ctx, mail, purpuse, now)
}

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
