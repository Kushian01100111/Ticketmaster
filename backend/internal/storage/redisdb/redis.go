package redisDB

import (
	"context"

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

	return rdb, nil
}
