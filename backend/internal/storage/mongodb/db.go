package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Kushian01100111/Tickermaster/internal/domain/booking"
	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/domain/ticket"
	"github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
)

func ConnectDB(dsn string, mongoDB string) (*mongo.Client, error) {
	conn, err := mongo.Connect(options.Client().
		ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := conn.Database(mongoDB)

	err = ensureCollections(ctx, db)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func updateCollections(ctx context.Context, db *mongo.Database) error {
	err := event.UpdateEventCollection(ctx, db)
	if err != nil {
		return err
	}

	err = ticket.UpdateTicketCollection(ctx, db)
	if err != nil {
		return err
	}

	err = booking.UpdateBookingCollection(ctx, db)
	if err != nil {
		return err
	}

	err = user.UpdateUserCollection(ctx, db)
	if err != nil {
		return err
	}

	return venue.UpdateVenueCollection(ctx, db)
}

func ensureCollections(ctx context.Context, db *mongo.Database) error {
	err := event.EnsureEventCollection(ctx, db)
	if err != nil {
		return err
	}

	err = ticket.EnsureTicketCollection(ctx, db)
	if err != nil {
		return err
	}

	err = booking.EnsureBookingCollection(ctx, db)
	if err != nil {
		return err
	}

	err = user.EnsureVenueCollection(ctx, db)
	if err != nil {
		return err
	}

	return venue.EnsureVenueCollection(ctx, db)
}
