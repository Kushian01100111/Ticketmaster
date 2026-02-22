package user

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	UserID bson.ObjectID `bson:"userID,omitempty"`

	UserName string `bson:"userName"`
	Email    string `bson:"email"`
	Password string `bson:"password,omitempty"`

	EasyLogin        bool      `bson:"easyLogin,omitempty"`
	FailedLoginCount int32     `bson:"failedLoginCount,omitempty"`
	LastFailedLogin  time.Time `bson:"lastFailedLogin,omitempty"`

	BookedEvents []bson.ObjectID `bson:"bookedEvents,omitempty"`
}
