package venue

import "go.mongodb.org/mongo-driver/bson/primitive"

type Venue struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	SeatType  string             `bson:"seatType"`
	SeatMapID primitive.ObjectID `bson:"venueId"`

	Address  string `bson:"address"`
	Capacity int32  `bson:"capacity"`
}
