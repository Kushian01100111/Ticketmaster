package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Kushian01100111/Tickermaster/internal/app/event"
	"github.com/Kushian01100111/Tickermaster/internal/app/venue"
	"github.com/Kushian01100111/Tickermaster/internal/config"
	http1 "github.com/Kushian01100111/Tickermaster/internal/http"
	"github.com/Kushian01100111/Tickermaster/internal/http/handlers"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"github.com/Kushian01100111/Tickermaster/internal/storage/mongodb"
)

func main() {
	var db *mongo.Client
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err = mongodb.ConnectDB(config.DSN, config.DB)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		_ = db.Disconnect(context.Background())
	}()

	addr := flag.String("addr", ":"+config.Port, "HTTP network address")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	eventRepo := repository.NewEventRepository(db.Database(config.DB), ctx)
	venueRepo := repository.NewVenueRepository(db.Database(config.DB), ctx)

	eventSvc := event.NewEventService(eventRepo, venueRepo)
	venueSvc := venue.NewVenueService(venueRepo)

	eventHandler := handlers.NewEventHandler(eventSvc)
	venueHandler := handlers.NewVenueHandler(venueSvc)

	r := http1.NewHandler(http1.RouterDep{
		EventDep: eventHandler,
		VenueDep: venueHandler},
		config)

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
