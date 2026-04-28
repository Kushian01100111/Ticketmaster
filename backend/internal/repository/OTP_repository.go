package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/otp"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidOTP        = errors.New("expired or invalid challenge")
	ErrChallengeNotFound = errors.New("challenge not found")
)

type OTPRepo interface {
	CreateOrReplace(ctx context.Context, ch otp.OTP) error
	GetActiveByEmail(ctx context.Context, mail, purpuse string) (*otp.OTP, error)
	IncAttempts(ctx context.Context, email, purpuse string) error
	Consume(ctx context.Context, email, purpuse string, when time.Time) error
}

type otpRepo struct {
	rdb *redis.Client
}

func NewOTPRepository(rdb *redis.Client) OTPRepo {
	return &otpRepo{rdb: rdb}
}

func (otp *otpRepo) CreateOrReplace(ctx context.Context, ch otp.OTP) error {
	/*
		filter := bson.M{
			"email":   ch.Email,
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
				"consumedAt": "",
			},
		}

		opts := options.UpdateOne().SetUpsert(true)

		_, err := otp.db.Collection("otp").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}

		return nil
	*/

	key := fmt.Sprintf("otp:%v:%v", ch.Email, ch.Purpuse)
	_, err := otp.rdb.JSONSet(ctx, key, "$", ch).Result()
	if err != nil {
		return err
	}

	_, err = otp.rdb.Expire(ctx, key, 10*time.Minute).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *otpRepo) GetActiveByEmail(ctx context.Context, mail string, purpuse string) (*otp.OTP, error) {
	/*
		filter := bson.M{
			"email":      mail,
			"purpuse":    purpuse,
			"expiresAt":  bson.M{"$gt": time.Now()},
			"consumedAt": bson.M{"$exists": false},
		}

		var otp *otp.OTPChallenge

		if err := r.db.Collection("otp").FindOne(ctx, filter).Decode(&otp); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrChallengeNotFound
			}
			return nil, err
		}

		return otp, nil
	*/
	var otp []otp.OTP
	key := fmt.Sprintf("otp:%v:%v", mail, purpuse)

	val, err := r.rdb.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &otp)
	if len(otp) > 0 {
		return &otp[0], nil
	}

	return nil, err
}

func (r *otpRepo) IncAttempts(ctx context.Context, email, purpuse string) error {
	/*
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
	*/

	key := fmt.Sprintf("otp:%v:%v", email, purpuse)
	_, err := r.rdb.JSONNumIncrBy(ctx, key, "$.attemps", 1).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *otpRepo) Consume(ctx context.Context, email, purpuse string, when time.Time) error {
	/*
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
	*/

	key := fmt.Sprintf("otp:%v:%v", email, purpuse)
	Exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil
	}

	if Exists == 1 {
		_, err := r.rdb.Del(ctx, key).Result()
		if err != nil {
			return nil
		}
	} else {
		return ErrInvalidOTP
	}

	return nil
}
