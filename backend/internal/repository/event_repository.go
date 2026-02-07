package repository

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventRepository interface {
	Create(event *event.Event, ctx context.Context) error
	Update(event *event.Event, ctx context.Context) error
	Delete(event *event.Event, ctx context.Context) error
	GetByID(event string, ctx context.Context) (*event.Event, error)      //Change name -> GetByID
	SearchByName(name string, ctx context.Context) ([]event.Event, error) //
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
func (s *mongoEventStorage) GetByID(event string, ctx context.Context) (*event.Event, error) {
	return nil, nil
}
func (s *mongoEventStorage) SearchByName(name string, ctx context.Context) ([]event.Event, error) {
	var res []event.Event
	return res, nil
}
