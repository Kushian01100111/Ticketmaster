package repository

import (
	"context"
	"errors"
	"regexp"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrDuplicateE       = errors.New("event is already created")
	ErrUpdateDuplicateE = errors.New("couldn't update event because other event has the same title or address")
	ErrEventNotFound    = errors.New("event couldn't be found")
)

type EventRepository interface {
	Create(event *event.Event, ctx context.Context) (bson.ObjectID, error)
	Update(event *event.Event, ctx context.Context) error
	Delete(id bson.ObjectID, ctx context.Context) error
	GetByID(id bson.ObjectID, ctx context.Context) (*event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
	SearchByParams(params *event.SearchEvent, ctx context.Context) ([]event.Event, error)
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

func (s *mongoEventStorage) SearchByParams(p *event.SearchEvent, ctx context.Context) ([]event.Event, error) {

	filter := bson.M{}

	if p.Currency != "" {
		filter["currency"] = p.Currency
	}

	if !p.VenueID.IsZero() {
		filter["venueId"] = p.VenueID
	}

	if p.Availability != "" {
		filter["availability"] = p.Availability
	}

	if !p.DateFrom.IsZero() || !p.DateTo.IsZero() {
		date := bson.M{}
		if !p.DateFrom.IsZero() {
			date["$gte"] = p.DateFrom
		}

		if !p.DateTo.IsZero() {
			date["$lte"] = p.DateTo
		}

		filter["startingDate"] = date
	}

	limit := 20

	sortField := "startingDate"
	if p.SortBy != "" {
		sortField = p.SortBy
	}

	order := int32(-1)
	if p.SortDir == 0 {
		order = 1
	}

	and := bson.A{}
	for _, t := range p.Tokens {
		t = regexp.QuoteMeta(t)
		and = append(and, bson.M{
			"$or": bson.A{
				bson.M{"title": bson.M{"$regex": t, "$options": "i"}},
				bson.M{"description": bson.M{"$regex": t, "$options": "i"}},
				bson.M{"venue.name": bson.M{"$regex": t, "$options": "i"}},
				bson.M{"venue.address": bson.M{"$regex": t, "$options": "i"}},
				bson.M{"performers": bson.M{"$elemMatch": bson.M{"$regex": t, "$options": "i"}}},
			},
		})
	}

	if len(and) > 0 {
		filter["$and"] = and
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: sortField, Value: order}, {Key: "_id", Value: order}})

	cur, err := s.db.Collection("event").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var res []event.Event
	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}
