package otp

import (
	"context"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/otp"
)

type OTPService interface {
	CreateOrReplace(ctx context.Context, ch OTPParams) error
	GetActiveByEmail(ctx context.Context, email string) (*otp.OTPChallange, error)
	IncAttempts(ctx context.Context, email string) error
	Consume(ctx context.Context, email string) error
}

type OTPParams struct {
	Email     string
	Purpuse   string
	CodeHash  string
	ExpiresAt time.Time
	Attempts  int32
	CreatedAt time.Time
}
