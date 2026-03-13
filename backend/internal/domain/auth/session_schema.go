package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func sessionSchema() bson.D {
	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"userId", "hash", "createdAt", "expiresAt"}},
			{Key: "properties", Value: bson.D{
				{Key: "userId", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "hash", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "createdAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
				{Key: "expiresAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
				{Key: "revokedAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
			}},
		},
	}}
}

func UpdateSessionCollecion(ctx context.Context, db *mongo.Database) error {
	validator := sessionSchema()

	cmd := bson.D{
		{Key: "collMod", Value: "session"},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return db.RunCommand(ctx, cmd).Err()
}

func EnsureSessionCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "session"}})
	if err != nil {
		return err
	}

	validator := sessionSchema()

	if len(existing) == 0 {
		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")
		if err := db.CreateCollection(ctx, "sesssion", opts); err != nil {
			return err
		}
	}

	coll := db.Collection("session")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "hast", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_hash_unique"),
		},
	})

	return err
}
