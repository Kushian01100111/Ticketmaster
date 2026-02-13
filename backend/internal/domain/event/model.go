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
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`

	StartingDate      time.Time `bson:"startingDate"`
	SalesStartingDate time.Time `bson:"salesStartingDate"`

	Currency  string `bson:"currency"`
	EventType string `bson:"eventType"`
	SeatType  string `bson:"seatType"`

	VenueID      bson.ObjectID `bson:"venueId"`
	Venue        Venue         `bson:"venue"`
	Performers   []string      `bson:"performers"`
	Status       string        `bson:"status"`
	Availability string        `bson:"availability"`
	Visibility   string        `bson:"visibility"`
}
