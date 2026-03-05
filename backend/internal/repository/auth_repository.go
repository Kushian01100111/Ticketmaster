package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type AuthRepository interface{}

type mongoAuthStorage struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &mongoAuthStorage{db: db}
}
