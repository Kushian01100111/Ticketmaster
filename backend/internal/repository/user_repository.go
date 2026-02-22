package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type UserRepository interface {
}

type mongoUserStorage struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserStorage{db: db}
}
