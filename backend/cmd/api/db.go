package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func connectDB(dsn string) (*mongo.Client, error) {
	db, err := mongo.Connect(options.Client().
		ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = db.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}
