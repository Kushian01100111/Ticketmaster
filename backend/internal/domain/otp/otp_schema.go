package otp

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func otpSchema() bson.D {
	purpuse := bson.D{
		{Key: "bsonType", Value: "string"},
		{Key: "enum", Value: bson.A{"login", "signUp"}},
	}

	return bson.D{{
		Key: "$jsonSchema",
		Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"email", "purpuse", "codeHash", "expiresAt", "attempts", "createdAt"}},
			{Key: "properties", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "bsonType", Value: "objectId"}}},
				{Key: "email", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "purpuse", Value: purpuse},
				{Key: "codeHash", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "expiresAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
				{Key: "attempts", Value: bson.D{{Key: "bsonType", Value: "int"}}},
				{Key: "consumedAt", Value: bson.D{{Key: "consumedAt", Value: "date"}}},
				{Key: "createdAt", Value: bson.D{{Key: "createdAt", Value: "date"}}},
			}},
		},
	}}
}

func UpdateOtpCollecion(ctx context.Context, db *mongo.Database) error {
	validator := otpSchema()

	cmd := bson.D{
		{Key: "collMod", Value: "otp"},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}

	return db.RunCommand(ctx, cmd).Err()
}

func EnsureOtpCollection(ctx context.Context, db *mongo.Database) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "otp"}})
	if err != nil {
		return err
	}

	validator := otpSchema()

	if len(existing) == 0 {
		opts := options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error").
			SetValidationLevel("strict")
		if err := db.CreateCollection(ctx, "otp", opts); err != nil {
			return err
		}
	} else {
		return UpdateOtpCollecion(ctx, db)
	}

	coll := db.Collection("session")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}, {Key: "purpuse", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_purpuse"),
		},
	})

	return err
}
