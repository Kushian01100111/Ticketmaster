package user

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func userSchema() bson.D {
	eventsSchema := bson.D{
		{Key: "bsonType", Value: "array"},
		{Key: "uniqueItems", Value: true},
		{Key: "items", Value: bson.D{
			{Key: "bsonType", Value: "objectId"},
		}},
	}

	roleSchema := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"costumer", "editor", "admin"}},
	}

	authMethods := bson.D{
		{Key: "bsonType", Value: "array"},
		{Key: "uniqueItems", Value: true},
		{Key: "items", Value: bson.D{
			{Key: "bsonType", Value: "string"},
			{Key: "enum", Value: bson.A{"password", "email_otp"}},
		}},
		{Key: "minItems", Value: 1},
	}

	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"email", "role", "authMethods"}},
			{Key: "properties", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "email", Value: bson.D{{Key: "bsonType", Value: "string"}}},

				{Key: "passwordHash", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "authMethods", Value: authMethods},

				{Key: "role", Value: roleSchema},

				{Key: "failedLoginCount", Value: bson.D{{Key: "bsonType", Value: "int"}}},
				{Key: "bookedEvents", Value: eventsSchema},
				{Key: "lastFailedLogin", Value: bson.D{{Key: "bsonType", Value: "date"}}},
			}},
		},
	}}
}

func UpdateUserCollection(ctx context.Context, db *mongo.Database) error {
	validator := userSchema()

	cmd := bson.D{
		{Key: "collMod", Value: "user"},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return db.RunCommand(ctx, cmd).Err()
}

func EnsureVenueCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "user"}})
	if err != nil {
		return err
	}

	validator := userSchema()

	if len(existing) == 0 {
		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")

		if err := db.CreateCollection(ctx, "user", opts); err != nil {
			return err
		}
	}

	coll := db.Collection("user")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
		},
	})

	return err
}
