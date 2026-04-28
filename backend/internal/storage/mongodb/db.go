package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Kushian01100111/Tickermaster/internal/domain/booking"
	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
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

	err = conn.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db := conn.Database(mongoDB)

	err = ensureCollections(context.Background(), db)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func UpdateCollections(ctx context.Context, db *mongo.Database) error {
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

	err = session.UpdateSessionCollecion(ctx, db)
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

	err = user.EnsureUserCollection(ctx, db)
	if err != nil {
		return err
	}

	err = session.EnsureSessionCollection(ctx, db)
	if err != nil {
		return err
	}

	return venue.EnsureVenueCollection(ctx, db)
}
