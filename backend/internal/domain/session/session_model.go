package session

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshSession struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"userId"`
	Hash      string        `bson:"hash"`
	ExpiresAt time.Time     `bson:"expiresAt"`
	CreatedAt time.Time     `bson:"createdAt"`
	RevokedAt *time.Time    `bson:"revokedAt,omitempty"`
}
