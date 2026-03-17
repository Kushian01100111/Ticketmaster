package repository

import (
	"context"
	"errors"

	"github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrDuplicateU   = errors.New("user is already created")
	ErrUserNotFound = errors.New("user couldn't be found")
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]user.User, error)
	Create(user *user.User, ctx context.Context) (bson.ObjectID, error)
	GetByID(id bson.ObjectID, ctx context.Context) (*user.User, error)
	GetByEmail(email string, ctx context.Context) (*user.User, error)
	UpdateUser(user *user.User, ctx context.Context) error
	DeleteUser(id bson.ObjectID, ctx context.Context) error

	FailedLogin(ctx context.Context, user *user.User) error
	ResetFailedLogin(ctx context.Context, user *user.User) error
}

type mongoUserStorage struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserStorage{db: db}
}

func (r *mongoUserStorage) GetAllUser(ctx context.Context) ([]user.User, error) {
	curr, err := r.db.Collection("user").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	var users []user.User
	if err := curr.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r mongoUserStorage) Create(user *user.User, ctx context.Context) (bson.ObjectID, error) {
	res, err := r.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.NilObjectID, ErrDuplicateU
		}
		return bson.NilObjectID, err
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, ErrPassingID
	}

	return id, nil
}

func (r mongoUserStorage) GetByID(id bson.ObjectID, ctx context.Context) (*user.User, error) {
	var out user.User
	err := r.db.Collection("user").
		FindOne(ctx, bson.D{{Key: "_id", Value: id}}).
		Decode(&out)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *mongoUserStorage) GetByEmail(email string, ctx context.Context) (*user.User, error) {
	var out user.User
	err := r.db.Collection("user").
		FindOne(ctx, bson.D{{Key: "email", Value: email}}).
		Decode(&out)

	if err != nil {
		return nil, err
	}

	return &out, err
}

func (r *mongoUserStorage) UpdateUser(user *user.User, ctx context.Context) error {
	filter := bson.M{"_id": user.ID}
	set := bson.M{
		"role":             user.Role,
		"authMethods":      user.AuthMethods,
		"failedLoginCount": user.FailedLoginCount,
		"lastFailedLogin":  user.LastFailedLogin,
		"bookedEvents":     user.BookedEvents,
	}

	if user.PasswordHash != nil {
		set["passwordHash"] = user.PasswordHash
	}

	update := bson.M{"$set": set}

	res, err := r.db.Collection("user").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *mongoUserStorage) DeleteUser(id bson.ObjectID, ctx context.Context) error {
	res, err := r.db.Collection("user").
		DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *mongoUserStorage) FailedLogin(ctx context.Context, user *user.User) error
func (r *mongoUserStorage) ResetFailedLogin(ctx context.Context, user *user.User) error
