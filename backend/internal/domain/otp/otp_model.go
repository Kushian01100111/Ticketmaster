package otp

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OTPChallenge struct {
	ID         bson.ObjectID `json:"_id,omitempty"`
	Email      string        `json:"email"`
	Purpuse    string        `json:"porpuse"`
	CodeHash   string        `json:"codeHash"`
	ExpiresAt  time.Time     `json:"expiresAt"`
	Attempts   int32         `json:"attempts"`
	ConsumedAt *time.Time    `json:"consumedAt,omitempty"`
	CreatedAt  time.Time     `json:"createdAt"`
}
