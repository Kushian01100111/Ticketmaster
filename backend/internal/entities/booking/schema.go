package booking

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func bookingSchema() bson.D {
	ticketSchema := bson.D{
		{Key: "bsonType", Value: "array"},
		{Key: "uniqueItems", Value: true},
		{Key: "items", Value: bson.D{
			{Key: "bsonType", Value: "objectId"},
			{Key: "minItems", Value: 1},
		}},
	}

	statusSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"reserved", "pendingPayment", "confirmed", "cancelled", "expired"}},
	}

	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"userId", "eventId", "total", "ticketList", "status", "dateCreated"}},
			{Key: "properties", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},

				{Key: "userId", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "eventId", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},

				{Key: "total", Value: bson.D{{Key: "bsonType", Value: "decimal"}}},

				{Key: "ticketList", Value: ticketSchema},
				{Key: "status", Value: statusSchema},

				{Key: "dateCreated", Value: bson.D{{Key: "bsonType", Value: "date"}}},
			}},
		},
	}}
}

func EnsureBookingCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "booking"}})
	if err != nil {
		return err
	}

	validator := bookingSchema()
	if len(existing) == 0 {
		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")

		if err := db.CreateCollection(ctx, "booking", opts); err != nil {
			return err
		}
	}

	coll := db.Collection("booking")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}},
			Options: options.Index().SetName("idx_relatedUser"),
		},
		{
			Keys:    bson.D{{Key: "eventId", Value: 1}},
			Options: options.Index().SetName("idx_relatedEvent"),
		},
		{
			Keys:    bson.D{{Key: "dateCreated", Value: 1}},
			Options: options.Index().SetName("idx_createdAt"),
		},
	})

	return err
}
