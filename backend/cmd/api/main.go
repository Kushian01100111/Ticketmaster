package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Kushian01100111/Tickermaster/internal/config"
	http1 "github.com/Kushian01100111/Tickermaster/internal/http"
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

	r := http1.NewHandler(config)

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
