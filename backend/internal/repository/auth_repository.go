package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RefreshSession struct {
	UserID    bson.ObjectID
	Hash      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, s RefreshSession) error
	RevokeRefreshToken(ctx context.Context, refreshTokenHash string) error
}

type mongoAuthStorage struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &mongoAuthStorage{db: db}
}

func (r *mongoAuthStorage) CreateRefreshToken(ctx context.Context, s RefreshSession) error
func (r *mongoAuthStorage) RevokeRefreshToken(ctx context.Context, refreshTokenHash string) error
