package repository

import (
	"context"
	"errors"

	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrDuplicate       = errors.New("venue is a already created")
	ErrUpdateDuplicate = errors.New("couldn't update venue because other venue has the same name or address")
	ErrPassingID       = errors.New("unexpected id type")
	ErrVenueNotFound   = errors.New("venue couldn't be found")
)

type VenueRepository interface {
	Create(venue *venue.Venue, ctx context.Context) (bson.ObjectID, error)
	Update(venue *venue.Venue, ctx context.Context) error
	Delete(id bson.ObjectID, ctx context.Context) error
	GetByID(id bson.ObjectID, ctx context.Context) (*venue.Venue, error)
	GetAll(ctx context.Context) ([]venue.Venue, error)
}

type mongoVenueStorage struct {
	db *mongo.Database
}

func NewVenueRepository(db *mongo.Database) VenueRepository {
	return &mongoVenueStorage{db: db}
}

func (s *mongoVenueStorage) Create(venue *venue.Venue, ctx context.Context) (bson.ObjectID, error) {
	res, err := s.db.Collection("venue").InsertOne(ctx, venue)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.NilObjectID, ErrDuplicate
		}
		return bson.NilObjectID, err
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, ErrPassingID
	}

	return id, err
}

func (s *mongoVenueStorage) Update(venue *venue.Venue, ctx context.Context) error {
	filter := bson.M{"_id": venue.ID}
	set := bson.M{
		"name":     venue.Name,
		"seatType": venue.SeatType,
		"address":  venue.Address,
		"capacity": venue.Capacity,
	}

	if !venue.SeatMapID.IsZero() {
		set["seatMapId"] = venue.SeatMapID
	}

	update := bson.M{"$set": set}

	res, err := s.db.Collection("venue").UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrUpdateDuplicate
		}
		return err
	}

	if res.MatchedCount == 0 {
		return ErrVenueNotFound
	}

	return nil
}

func (s *mongoVenueStorage) Delete(id bson.ObjectID, ctx context.Context) error {
	res, err := s.db.Collection("venue").
		DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrVenueNotFound
	}

	return nil
}

func (s *mongoVenueStorage) GetByID(id bson.ObjectID, ctx context.Context) (*venue.Venue, error) {
	var out venue.Venue
	err := s.db.Collection("venue").
		FindOne(ctx, bson.D{{Key: "_id", Value: id}}).
		Decode(&out)

	if err != nil {
		return nil, err
	}

	return &out, err
}

func (s *mongoVenueStorage) GetAll(ctx context.Context) ([]venue.Venue, error) {
	cur, err := s.db.Collection("venue").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var venues []venue.Venue
	if err := cur.All(ctx, &venues); err != nil {
		return nil, err
	}
	return venues, err
}
