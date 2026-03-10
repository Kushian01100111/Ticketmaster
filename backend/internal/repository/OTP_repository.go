package repository

import (
	"context"
	"time"
)

type OTPChallange struct {
	Email      string
	CodeHash   string
	ExpiresAt  time.Time
	Attempts   int32
	ConsumedAt *time.Time
	CreatedAt  time.Time
}

type OTPRepo interface {
	CreateOrReplace(ctx context.Context, ch OTPChallange) error
	GetActiveByEmail(ctx context.Context, email string) (*OTPChallange, error)
	IncAttempts(ctx context.Context, email string) error
	Consume(ctx context.Context, email string, when time.Time) error
}

type otpRepo struct{}

func NewOTPRepository() OTPRepo {
	return otpRepo{}
}

func (otp *otpRepo) CreateOrReplace(ctx context.Context, ch OTPChallange) error
func (otp *otpRepo) GetActiveByEmail(ctx context.Context, email string) (*OTPChallange, error)
func (otp *otpRepo) IncAttempts(ctx context.Context, email string) error
func (otp *otpRepo) Consume(ctx context.Context, email string, when time.Time) error
