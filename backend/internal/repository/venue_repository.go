package repository

import (
	"context"
	"errors"

	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrDuplicate     = errors.New("venue is a already created")
	ErrPassingID     = errors.New("unexpected id type")
	ErrVenueNotFound = errors.New("venue could'nt be found")
)

type VenueRepository interface {
	Create(venue *venue.Venue, ctx context.Context) (bson.ObjectID, error)
	Update(venue *venue.Venue, ctx context.Context) error
	Delete(venue *venue.Venue, ctx context.Context) error
	GetByID(idHex string, ctx context.Context) (*venue.Venue, error)
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
	update := bson.M{"$set": bson.M{
		"name":     venue.Name,
		"seatType": venue.SeatType,
		"address":  venue.Address,
		"capacity": venue.Capacity,
	}}

	if !venue.SeatMapID.IsZero() {
		update = bson.M{"$set": bson.M{
			"name":      venue.ID,
			"seatType":  venue.SeatType,
			"seatMapId": venue.SeatMapID,
			"address":   venue.Address,
			"capacity":  venue.Capacity,
		}}

	}

	res, err := s.db.Collection("venue").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrVenueNotFound
	}

	return nil
}

func (s *mongoVenueStorage) Delete(venue *venue.Venue, ctx context.Context) error {
	return nil
}

func (s *mongoVenueStorage) GetByID(idHex string, ctx context.Context) (*venue.Venue, error) {
	id, err := bson.ObjectIDFromHex(idHex)
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
