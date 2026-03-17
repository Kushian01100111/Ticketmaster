package repository

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, s auth.RefreshSession) error
	GetByHash(ctx context.Context, hash string) (*auth.RefreshSession, error)
	RevokeRefreshToken(ctx context.Context, refreshTokenHash string) error
	RevokeAllByUserID(ctx context.Context, oid bson.ObjectID) error
}

type mongoAuthStorage struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &mongoAuthStorage{db: db}
}

func (r *mongoAuthStorage) CreateRefreshToken(ctx context.Context, s auth.RefreshSession) error
func (r *mongoAuthStorage) GetByHash(ctx context.Context, hash string) (*auth.RefreshSession, error)
func (r *mongoAuthStorage) RevokeRefreshToken(ctx context.Context, refreshTokenHash string) error
func (r *mongoAuthStorage) RevokeAllByUserID(ctx context.Context, oid bson.ObjectID) error
