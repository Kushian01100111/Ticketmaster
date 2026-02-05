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
	GetByName(event string) (*event.Event, error)     //Change name -> GetByID
	SearchByName(event string) ([]event.Event, error) //
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
func (s *mongoEventStorage) Update(event *event.Event) error
func (s *mongoEventStorage) Delete(event *event.Event) error
func (s *mongoEventStorage) GetByName(event string) (*event.Event, error)
func (s *mongoEventStorage) SearchByName(event string) ([]event.Event, error)
