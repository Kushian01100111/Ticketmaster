package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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
