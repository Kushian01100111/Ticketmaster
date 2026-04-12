package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/otp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrInvalidOTP        = errors.New("invalid challenge")
	ErrChallangeNotFound = errors.New("challenge not found")
)

type OTPRepo interface {
	CreateOrReplace(ctx context.Context, ch otp.OTPChallange) error
	GetActiveByEmail(ctx context.Context, mail string, purpuse string) (*otp.OTPChallange, error)
	IncAttempts(ctx context.Context, email string) error
	Consume(ctx context.Context, email string, when time.Time) error
}

type otpRepo struct {
	db *mongo.Database
}

func NewOTPRepository(db *mongo.Database) OTPRepo {
	return &otpRepo{db: db}
}

func (otp *otpRepo) CreateOrReplace(ctx context.Context, ch otp.OTPChallange) error {
	filter := bson.M{
		"email":   strings.ToLower(strings.TrimSpace(ch.Email)),
		"purpuse": ch.Purpuse,
	}

	update := bson.M{
		"$set": bson.M{
			"email":     filter["email"],
			"purpuse":   ch.Purpuse,
			"codeHash":  ch.CodeHash,
			"expiresAt": ch.ExpiresAt,
			"attempts":  ch.Attempts,
			"createdAt": ch.CreatedAt,
		},
		"$unset": bson.M{
			"cosumeAt": "",
		},
	}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := otp.db.Collection("otp").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *otpRepo) GetActiveByEmail(ctx context.Context, mail string, purpuse string) (*otp.OTPChallange, error) {
	filter := bson.M{
		"email":      mail,
		"purpuse":    purpuse,
		"expiresAt":  bson.M{"$gt": time.Now()},
		"consumedAt": bson.M{"$exists": false},
	}

	var otp *otp.OTPChallange

	if err := r.db.Collection("otp").FindOne(ctx, filter).Decode(&otp); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrChallangeNotFound
		}
		return nil, err
	}

	return otp, nil
}

func (r *otpRepo) IncAttempts(ctx context.Context, email string) error {
	filter := bson.M{
		"email":      email,
		"consumedAt": bson.M{"$exists": false},
		"expiresAt":  bson.M{"$gt": time.Now()},
	}

	update := bson.M{
		"$inc": bson.M{"attempts": 1},
	}

	res, err := r.db.Collection("otp").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrInvalidOTP
	}

	return nil
}

func (r *otpRepo) Consume(ctx context.Context, email string, when time.Time) error {
	filter := bson.M{
		"email":      email,
		"consumedAt": bson.M{"$exists": false},
		"expiresAt":  bson.M{"$gt": time.Now()},
	}

	update := bson.M{
		"$set": bson.M{"consumedAt": when},
	}

	res, err := r.db.Collection("otp").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrInvalidOTP
	}

	return nil
}
