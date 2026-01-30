package event

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Venue struct {
	Name     string `bson:"name"`
	Address  string `bson:"address"`
	Capacity int32  `bson:"int"`
}

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`

	Date       time.Time `bson:"startAt,omitempty"`
	SalesStart time.Time `bson:"salesStarttAt,omitempty"`

	Currency  string `bson:"currency"`
	EventType string `bson:"eventType,omitempty"`
	SeatType  string `bson:"seatType"`

	VenueID      primitive.ObjectID `bson:"venueId,omitempty"`
	Venue        Venue              `bson:"venue,omitempty"`
	Performers   []string           `bson:"artists,omitempty"`
	Status       string             `bson:"status,omitempty"`
	Availability string             `bson:"availability,omitempty"`
	Visibility   string             `bson:"visibility,omitempty"`
}
