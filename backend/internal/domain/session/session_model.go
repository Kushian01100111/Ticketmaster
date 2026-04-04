package session

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshSession struct {
	UserID    bson.ObjectID `bson:"userId,omitempty"`
	Hash      string        `bson:"hash"`
	ExpiresAt time.Time     `bson:"expiresAt"`
	CreatedAt time.Time     `bson:"createdAt"`
	RevokedAt *time.Time    `bson:"revokedAt"`
}
