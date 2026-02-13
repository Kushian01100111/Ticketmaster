package repository

import (
	"context"
	"errors"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrDuplicateE       = errors.New("event is a already created")
	ErrUpdateDuplicateE = errors.New("couldn't update event because other event has the same title or address")
	ErrEventNotFound    = errors.New("event couldn't be found")
)

type EventRepository interface {
	Create(event *event.Event, ctx context.Context) (bson.ObjectID, error)
	Update(event *event.Event, ctx context.Context) error
	Delete(id bson.ObjectID, ctx context.Context) error
	GetByID(id bson.ObjectID, ctx context.Context) (*event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
	SearchByName(name string, ctx context.Context) ([]event.Event, error)
}

type mongoEventStorage struct {
	db *mongo.Database
}

func NewEventRepository(db *mongo.Database) EventRepository {
	return &mongoEventStorage{db: db}
}

func (s *mongoEventStorage) Create(event *event.Event, ctx context.Context) (bson.ObjectID, error) {
	res, err := s.db.Collection("event").InsertOne(ctx, event)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.NilObjectID, ErrDuplicateE
		}
		return bson.NilObjectID, err
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, ErrPassingID
	}

	return id, nil
}
func (s *mongoEventStorage) Update(event *event.Event, ctx context.Context) error {
	filter := bson.M{"_id": event.ID}
	set := bson.M{
		"title":             event.Title,
		"description":       event.Description,
		"startingDate":      event.StartingDate,
		"salesStartingDate": event.SalesStartingDate,
		"currency":          event.Currency,
		"eventType":         event.EventType,
		"SeatType":          event.SeatType,
		"venueId":           event.VenueID,
		"venue": bson.M{
			"name":     event.Venue.Name,
			"address":  event.Venue.Address,
			"capacity": event.Venue.Capacity,
		},
		"performers":   event.Performers,
		"status":       event.Status,
		"availability": event.Availability,
		"Visibility":   event.Visibility,
	}

	update := bson.M{"$set": set}

	res, err := s.db.Collection("event").UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateE
		}
		return err
	}

	if res.MatchedCount == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *mongoEventStorage) Delete(id bson.ObjectID, ctx context.Context) error {
	res, err := s.db.Collection("event").
		DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrEventNotFound
	}

	return nil
}
func (s *mongoEventStorage) GetByID(id bson.ObjectID, ctx context.Context) (*event.Event, error) {
	var out event.Event
	err := s.db.Collection("event").
		FindOne(ctx, bson.D{{Key: "_id", Value: id}}).
		Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (s *mongoEventStorage) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	curr, err := s.db.Collection("event").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	var events []event.Event
	if err := curr.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, err
}

func (s *mongoEventStorage) SearchByName(name string, ctx context.Context) ([]event.Event, error) {
	var res []event.Event
	return res, nil
}
