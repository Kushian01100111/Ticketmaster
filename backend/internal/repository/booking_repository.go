package repository

import "github.com/redis/go-redis/v9"

type BookingRepo interface {
}

type bookingRepo struct {
	rdb *redis.Client
}

func NewBookingRepo(rdb *redis.Client) BookingRepo {
	return &bookingRepo{rdb: rdb}
}
