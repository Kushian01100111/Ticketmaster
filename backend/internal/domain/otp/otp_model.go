package otp

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OTPChallange struct {
	ID         bson.ObjectID `json:"_id,omitempty"`
	Email      string        `json:"email"`
	Porpuse    string        `json:"porpuse"`
	CodeHash   string        `json:"codeHash"`
	ExpiresAt  time.Time     `json:"expiresAt"`
	Attempts   int32         `json:"attempts"`
	ConsumedAt *time.Time    `json:"consumeAt,omitempty"`
	CreatedAt  time.Time     `json:"createdAt"`
}
