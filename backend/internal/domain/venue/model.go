package venue

import "go.mongodb.org/mongo-driver/v2/bson"

type Venue struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`

	SeatType  string        `bson:"seatType"`
	SeatMapID bson.ObjectID `bson:"seatMapId,omitempty"`

	Address  string `bson:"address"`
	Capacity int32  `bson:"capacity"`
}
