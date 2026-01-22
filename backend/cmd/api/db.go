package main

import (
	"context"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/entities/venue"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func connectDB(dsn string, mongoDB string) (*mongo.Client, error) {
	conn, err := mongo.Connect(options.Client().
		ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = conn.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := conn.Database(mongoDB)

	err = EnsureCollections(ctx, db)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func EnsureCollections(ctx context.Context, db *mongo.Database) error {
	return venue.EnsureVenueCollection(ctx, db)
}
