package otp

import (
	"context"

	"github.com/redis/go-redis/v9"
)

/*
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
				{Key: "consumedAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
				{Key: "createdAt", Value: bson.D{{Key: "bsonType", Value: "date"}}},
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

	coll := db.Collection("otp")
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}, {Key: "purpuse", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_purpuse"),
		},
	})

	return err
}

*/

func EnsureOTPRedis(ctx context.Context, rdb *redis.Client) error {
	indeces, err := rdb.FT_List(ctx).Result()
	if err != nil {
		return err
	}

	exists := false
	for _, idx := range indeces {
		if idx == "idx:otp" {
			exists = true
			break
		}
	}

	if !exists {
		_, err := rdb.FTCreate(
			ctx,
			"idx:otp",
			&redis.FTCreateOptions{
				OnJSON: true,
				Prefix: []interface{}{"otp"},
			},

			&redis.FieldSchema{
				FieldName: "$.email",
				As:        "email",
				FieldType: redis.SearchFieldTypeTag,
			},

			&redis.FieldSchema{
				FieldName: "$.purpuse",
				As:        "purpuse",
				FieldType: redis.SearchFieldTypeText,
			},

			&redis.FieldSchema{
				FieldName: "$.expiresAt",
				As:        "expiresAt",
				FieldType: redis.SearchFieldTypeNumeric,
			},
		).Result()
		if err != nil {
			return err
		}
	}

	return nil
}
