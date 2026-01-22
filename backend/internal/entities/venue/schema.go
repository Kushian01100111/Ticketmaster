package venue

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureVenueCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "venue"}})
	if err != nil {
		return err
	}

	if len(existing) == 0 {
		validator := bson.D{{
			Key: "$jsonSchema",
			Value: bson.D{
				{Key: "bsonType", Value: "object"},
				{Key: "required", Value: bson.A{"address", "capacity"}},
				{Key: "properties", Value: bson.D{
					{Key: "address", Value: bson.D{{Key: "bsonType", Value: "string"}}},
					{Key: "capacity", Value: bson.D{{Key: "bsonType", Value: "int32"}}},
				}},
			},
		}}

		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")

		if err := db.CreateCollection(ctx, "venue", opts); err != nil {
			return err
		}
	}

	coll := db.Collection("venue")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "address", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "capacity", Value: 1}},
		},
	})

	return err
}
