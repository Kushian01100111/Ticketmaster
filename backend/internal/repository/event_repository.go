package repository

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventRepository interface {
	Create(event *event.Event, ctx context.Context) error
	Update(event *event.Event, ctx context.Context) error
	Delete(event *event.Event, ctx context.Context) error
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

func (s *mongoEventStorage) Create(event *event.Event, ctx context.Context) error {
	_, err := s.db.Collection("event").InsertOne(ctx, event)
	return err
}
func (s *mongoEventStorage) Update(event *event.Event, ctx context.Context) error {
	return nil
}

func (s *mongoEventStorage) Delete(event *event.Event, ctx context.Context) error {
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
