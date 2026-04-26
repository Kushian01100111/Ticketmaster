package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Kushian01100111/Tickermaster/internal/app/auth"
	"github.com/Kushian01100111/Tickermaster/internal/app/email"
	"github.com/Kushian01100111/Tickermaster/internal/app/event"
	otpChallenge "github.com/Kushian01100111/Tickermaster/internal/app/otpChallange"
	"github.com/Kushian01100111/Tickermaster/internal/app/user"
	"github.com/Kushian01100111/Tickermaster/internal/app/venue"
	"github.com/Kushian01100111/Tickermaster/internal/config"
	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
	http1 "github.com/Kushian01100111/Tickermaster/internal/http"
	"github.com/Kushian01100111/Tickermaster/internal/http/handlers"
	"github.com/Kushian01100111/Tickermaster/internal/http/middleware"
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

	JWTManager, err := session.NewJWTManager(session.JWTConfig{
		Secret:    config.JWTSECRET,
		Issuer:    "booking",
		Audience:  "booking-web",
		AccessTTL: 15 * time.Minute,
		ClockSkew: 30 * time.Minute,
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Redis -> planing
	hasher := session.NewBcryptHasher(0)
	middleware := middleware.NewAuthMiddleware(JWTManager)

	authRepo := repository.NewAuthRepository(db.Database(config.DB))
	eventRepo := repository.NewEventRepository(db.Database(config.DB))
	venueRepo := repository.NewVenueRepository(db.Database(config.DB))
	userRepo := repository.NewUserRepository(db.Database(config.DB))
	otpRepo := repository.NewOTPRepository(db.Database(config.DB))
	emailRepo := email.NewEmailSender(config.ResendAPIKey, config.EmailFrom)

	eventSvc := event.NewEventService(eventRepo, venueRepo)
	venueSvc := venue.NewVenueService(venueRepo)
	userSvc := user.NewUserService(userRepo)
	otpSrv := otpChallenge.NewOTPService(otpRepo, userRepo)
	authSrv := auth.NewAuthService(
		otpSrv,
		authRepo,
		userSvc,
		emailRepo,
		hasher,
		JWTManager,
		auth.AuthConfig{OTPTTL: 10 * time.Minute, RefreshTTL: 30 * 24 * time.Hour},
	)

	authHandler := handlers.NewAuthHandler(authSrv)
	eventHandler := handlers.NewEventHandler(eventSvc)
	venueHandler := handlers.NewVenueHandler(venueSvc)
	userHandler := handlers.NewUserHandler(userSvc)

	r := http1.NewHandler(http1.RouterDep{
		AuthHandler: authHandler,
		EventDep:    eventHandler,
		VenueDep:    venueHandler,
		UserDep:     userHandler},
		config,
		middleware)

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
