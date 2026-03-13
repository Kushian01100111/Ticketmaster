package user

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Email string `bson:"email"`
	Role  string `bson:"role"`

	PasswordHash *string  `bson:"passwordHash,omitempty"`
	AuthMethods  []string `bson:"authMethods"`

	FailedLoginCount int32      `bson:"failedLoginCount,omitempty"`
	LastFailedLogin  *time.Time `bson:"lastFailedLogin,omitempty"`

	BookedEvents []bson.ObjectID `bson:"bookedEvents,omitempty"`
}
