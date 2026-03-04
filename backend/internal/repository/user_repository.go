package repository

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]user.User, error)
	Create(user *user.User, ctx context.Context) (bson.ObjectID, error)
	GetByID(id bson.ObjectID, ctx context.Context) (*user.User, error)
	UpdateUser(user *user.User, ctx context.Context) error
	DeleteUser(id bson.ObjectID, ctx context.Context) error
}

type mongoUserStorage struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserStorage{db: db}
}

func (r mongoUserStorage) GetAllUser(ctx context.Context) ([]user.User, error) {

}

func (r mongoUserStorage) Create(user *user.User, ctx context.Context) (bson.ObjectID, error) {

}

func (r mongoUserStorage) GetByID(id bson.ObjectID, ctx context.Context) (*user.User, error) {

}

func (r mongoUserStorage) UpdateUser(user *user.User, ctx context.Context) error {

}

func (r mongoUserStorage) DeleteUser(id bson.ObjectID, ctx context.Context) error {

}
