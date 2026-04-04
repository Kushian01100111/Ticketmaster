package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrDuplicateS     = errors.New("session is already created")
	ErrSessioNotFound = errors.New("session was not found")
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, s session.RefreshSession) error
	GetByHash(ctx context.Context, hash string) (*session.RefreshSession, error)
	RevokeRefreshToken(ctx context.Context, refreshTokenHash string) error
	RevokeAllByUserID(ctx context.Context, oid bson.ObjectID) error
}

type mongoAuthStorage struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &mongoAuthStorage{db: db}
}

func (r *mongoAuthStorage) CreateRefreshToken(ctx context.Context, s session.RefreshSession) error {
	res, err := r.db.Collection("session").InsertOne(ctx, s)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateS
		}
		return err
	}

	_, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return ErrPassingID
	}
	return nil
}
func (r *mongoAuthStorage) GetByHash(ctx context.Context, hash string) (*session.RefreshSession, error) {
	var out session.RefreshSession
	err := r.db.Collection("session").
		FindOne(ctx, bson.D{{Key: "hash", Value: hash}}).
		Decode(out)

	if err != nil {
		return nil, err
	}
	return &out, nil
}
func (r *mongoAuthStorage) RevokeRefreshToken(ctx context.Context, refreshTokenHash string) error {
	filter := bson.M{"hash": refreshTokenHash}
	set := bson.M{
		"revokedAt": time.Now(),
	}
	update := bson.M{"$set": set}

	res, err := r.db.Collection("session").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrSessioNotFound
	}

	return nil
}
func (r *mongoAuthStorage) RevokeAllByUserID(ctx context.Context, oid bson.ObjectID) error {
	when := time.Now()

	filter := bson.M{
		"userId":    oid,
		"revokedAt": bson.M{"$exists": false},
	}

	update := bson.M{"$set": bson.M{"revokedAt": when}}

	res, err := r.db.Collection("session").UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return ErrSessioNotFound
	}

	return nil
}
