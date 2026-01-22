package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	var db *mongo.Client
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err = connectDB(config.DSN, config.DB)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		_ = db.Disconnect(context.Background())
	}()

	addr := flag.String("addr", ":"+config.Port, "HTTP network address")

	r := NewHandler(config)

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
