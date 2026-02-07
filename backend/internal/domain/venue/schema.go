package venue

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func venueSchema() bson.D {
	seatTypeSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"seated", "standing", "mixed"}},
	}

	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"name", "seatType", "address", "capacity"}},
			{Key: "properties", Value: bson.D{
				{Key: "name", Value: bson.D{{Key: "bsonType", Value: "string"}}},

				{Key: "seatType", Value: seatTypeSchema},
				{Key: "seatMapID", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},

				{Key: "address", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "capacity", Value: bson.D{{Key: "bsonType", Value: "int"}}},
			}},
		},
	}}
}

func UpdateVenueCollection(ctx context.Context, db *mongo.Database) error {
	validator := venueSchema()

	cmd := bson.D{
		{Key: "collMod", Value: "user"},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return db.RunCommand(ctx, cmd).Err()
}

func EnsureVenueCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "venue"}})
	if err != nil {
		return err
	}

	validator := venueSchema()
	if len(existing) == 0 {
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
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
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
