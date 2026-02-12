package event

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Venue struct {
	Name     string `bson:"name"`
	Address  string `bson:"address"`
	Capacity int32  `bson:"int"`
}

type Event struct {
	ID          bson.ObjectID `bson:"_id"`
	Title       string        `bson:"name"`
	Description string        `bson:"description"`

	StartingDate time.Time `bson:"startingDate"`
	SalesStart   time.Time `bson:"salesStarttAt"`

	Currency  string `bson:"currency"`
	EventType string `bson:"eventType"`
	SeatType  string `bson:"seatType"`

	VenueID      bson.ObjectID `bson:"venueId"`
	Venue        Venue         `bson:"venue,omitempty"`
	Performers   []string      `bson:"performers"`
	Status       string        `bson:"status"`
	Availability string        `bson:"availability,"`
	Visibility   string        `bson:"visibility"`
}
