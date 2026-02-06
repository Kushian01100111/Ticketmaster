package repository

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventRepository interface {
	Create(event *event.Event) error
	Update(event *event.Event) error
	Delete(event *event.Event) error
	GetByID(event string) (*event.Event, error)      //Change name -> GetByID
	SearchByName(name string) ([]event.Event, error) //
}

type mongoEventStorage struct {
	db  *mongo.Database
	ctx context.Context
}

func NewEventRepository(db *mongo.Database, ctx context.Context) EventRepository {
	return &mongoEventStorage{db: db, ctx: ctx}
}

func (s *mongoEventStorage) Create(event *event.Event) error {
	_, err := s.db.Collection("event").InsertOne(s.ctx, event)
	return err
}
func (s *mongoEventStorage) Update(event *event.Event) error {
	return nil
}

func (s *mongoEventStorage) Delete(event *event.Event) error {
	return nil
}
func (s *mongoEventStorage) GetByID(event string) (*event.Event, error) {
	return nil, nil
}
func (s *mongoEventStorage) SearchByName(name string) ([]event.Event, error) {
	var res []event.Event
	return res, nil
}
