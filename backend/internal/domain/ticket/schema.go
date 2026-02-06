package ticket

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ticketSchema() bson.D {
	statusSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"hold", "available", "sold"}},
	}

	seatSchema := bson.D{
		{Key: "bsonType", Value: "object"},
		{Key: "required", Value: bson.A{"section", "row", "number"}},
		{Key: "properties", Value: bson.D{
			{Key: "section", Value: bson.D{{Key: "bsonType", Value: "string"}}},
			{Key: "row", Value: bson.D{{Key: "bsonType", Value: "number"}}},
			{Key: "number", Value: bson.D{{Key: "bsonType", Value: "int"}}},
		}},
	}

	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"userId", "eventId", "price", "status", "seatData"}},
			{Key: "properties", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},

				{Key: "userId", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "eventId", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},

				{Key: "price", Value: bson.D{{Key: "bsonType", Value: "decimal"}}},
				{Key: "status", Value: statusSchema},
				{Key: "seatData", Value: seatSchema},
			}},
		},
	}}
}

func UpdateTicketCollection(ctx context.Context, db *mongo.Database) error {
	validator := ticketSchema()

	cmd := bson.D{
		{Key: "collMod", Value: "ticket"},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return db.RunCommand(ctx, cmd).Err()
}

func EnsureTicketCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "ticket"}})
	if err != nil {
		return err
	}

	validator := ticketSchema()

	if len(existing) == 0 {
		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")

		if err := db.CreateCollection(ctx, "ticket", opts); err != nil {
			return err
		}
	}

	coll := db.Collection("ticket")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
	})

	return err
}
