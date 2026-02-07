package repository

import (
	"context"
	"errors"

	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrDuplicate = errors.New("venue is a already created")
	ErrPassingID = errors.New("unexpected id type")
)

type VenueRepository interface {
	Create(venue *venue.Venue, ctx context.Context) (primitive.ObjectID, error)
	Update(venue *venue.Venue, ctx context.Context) error
	Delete(venue *venue.Venue, ctx context.Context) error
	GetByID(idHex string, ctx context.Context) (*venue.Venue, error)
}

type mongoVenueStorage struct {
	db *mongo.Database
}

func NewVenueRepository(db *mongo.Database) VenueRepository {
	return &mongoVenueStorage{db: db}
}

func (s *mongoVenueStorage) Create(venue *venue.Venue, ctx context.Context) (primitive.ObjectID, error) {
	res, err := s.db.Collection("venue").InsertOne(ctx, venue)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, ErrDuplicate
		}
		return primitive.NilObjectID, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, ErrPassingID
	}

	return id, err
}

func (s *mongoVenueStorage) Update(venue *venue.Venue, ctx context.Context) error {
	return nil
}

func (s *mongoVenueStorage) Delete(venue *venue.Venue, ctx context.Context) error {
	return nil
}

func (s *mongoVenueStorage) GetByID(idHex string, ctx context.Context) (*venue.Venue, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	var out venue.Venue
	err = s.db.Collection("venue").
		FindOne(ctx, bson.D{{Key: "_id", Value: id}}).
		Decode(&out)

	if err != nil {
		return nil, err
	}

	return &out, err
}
