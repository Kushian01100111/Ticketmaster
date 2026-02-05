package repository

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type VenueRepository interface {
	Create(venue *venue.Venue) error
	Update(venue *venue.Venue) error
	Delete(venue *venue.Venue) error
	GetByID(idHex string) (*venue.Venue, error)
}

type mongoVenueStorage struct {
	db  *mongo.Database
	ctx context.Context
}

func NewVenueRepository(db *mongo.Database, ctx context.Context) VenueRepository {
	return &mongoVenueStorage{db: db, ctx: ctx}
}

func (s *mongoVenueStorage) Create(venue *venue.Venue) error
func (s *mongoVenueStorage) Update(venue *venue.Venue) error
func (s *mongoVenueStorage) Delete(venue *venue.Venue) error
func (s *mongoVenueStorage) GetByID(idHex string) (*venue.Venue, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	var out venue.Venue
	err = s.db.Collection("venue").
		FindOne(s.ctx, bson.D{{Key: "_id", Value: id}}).
		Decode(&out)

	if err != nil {
		return nil, err
	}

	return &out, err
}
