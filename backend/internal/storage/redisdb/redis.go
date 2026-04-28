package redisDB

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/otp"
	"github.com/redis/go-redis/v9"
)

type RedisSecrets struct {
	Addr     string
	User     string
	Password string
}

func ConnectRDB(secrets RedisSecrets) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     secrets.Addr,
		Username: secrets.User,
		Password: secrets.Password,
		DB:       0,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	err = ensureSchemas(ctx, rdb)
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func ensureSchemas(ctx context.Context, rdb *redis.Client) error {
	err := otp.EnsureOTPRedis(ctx, rdb)
	if err != nil {
		return err
	}

	return nil
}
