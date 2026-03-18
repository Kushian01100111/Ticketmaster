package repository

import (
	"context"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/otp"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OTPRepo interface {
	CreateOrReplace(ctx context.Context, ch otp.OTPChallange) (string, error)
	GetActiveByEmail(ctx context.Context, email string) (*otp.OTPChallange, error)
	IncAttempts(ctx context.Context, email string) error
	Consume(ctx context.Context, email string, when time.Time) error
}

type otpRepo struct {
	db *mongo.Database
}

func NewOTPRepository(db *mongo.Database) OTPRepo {
	return &otpRepo{db: db}
}

func (otp *otpRepo) CreateOrReplace(ctx context.Context, ch otp.OTPChallange) (string, error)

func (otp *otpRepo) GetActiveByEmail(ctx context.Context, email string) (*otp.OTPChallange, error)
func (otp *otpRepo) IncAttempts(ctx context.Context, email string) error
func (otp *otpRepo) Consume(ctx context.Context, email string, when time.Time) error
