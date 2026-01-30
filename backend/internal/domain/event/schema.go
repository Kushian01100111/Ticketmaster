package event

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func eventSchema() bson.D {
	artistSchema := bson.D{
		{Key: "bsonType", Value: "array"},
		{Key: "uniqueItems", Value: true},
		{Key: "items", Value: bson.D{
			{Key: "bsonType", Value: "string"},
			{Key: "minItems", Value: 1},
		}},
	}

	eventTypeSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"concert", "recital", "solo recitals", "operatic productions"}},
	}

	statusSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"draft", "published", "cancelled", "postpond"}},
	}

	visibilitySchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"public", "unlisted", "private"}},
	}

	availabilitySchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"available", "soldOut"}},
	}

	venueShema := bson.D{
		{Key: "bsonType", Value: "object"},
		{Key: "Properties", Value: bson.D{
			{Key: "name", Value: "string"},
			{Key: "address", Value: "string"},
			{Key: "capacity", Value: "int"},
		}},
	}

	seatTypeSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"seated", "standing", "mixed"}},
	}

	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"name", "startAt", "description", "eventType", "venue", "status", "visibility", "salesStarttAt", "seatType"}},
			{Key: "properties", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "name", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "description", Value: bson.D{{Key: "bsonType", Value: "string"}}},

				{Key: "startAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
				{Key: "salesStarttAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},

				{Key: "currency", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "eventType", Value: eventTypeSchema},
				{Key: "seatType", Value: seatTypeSchema},
				// Missing: seatMap, pricing and pricingMap

				{Key: "venueId", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "venue", Value: venueShema},
				{Key: "artists", Value: artistSchema},
				{Key: "status", Value: statusSchema},
				{Key: "Availability", Value: availabilitySchema},
				{Key: "visibility", Value: visibilitySchema},
			}},
		},
	}}
}

func UpdateEventCollection(ctx context.Context, db *mongo.Database) error {
	validator := eventSchema()

	cmd := bson.D{
		{Key: "collMod", Value: "event"},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return db.RunCommand(ctx, cmd).Err()
}

func EnsureEventCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "event"}})
	if err != nil {
		return err
	}

	validator := eventSchema()
	if len(existing) == 0 {
		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")

		if err := db.CreateCollection(ctx, "event", opts); err != nil {
			return err
		}
	}

	coll := db.Collection("event")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "date", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_eventName"),
		},
		{
			Keys:    bson.D{{Key: "eventType", Value: 1}},
			Options: options.Index().SetName("idx_eventType"),
		},
	})

	return err
}
