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
	"github.com/Kushian01100111/Tickermaster/internal/app/booking"
	"github.com/Kushian01100111/Tickermaster/internal/app/email"
	"github.com/Kushian01100111/Tickermaster/internal/app/event"
	"github.com/Kushian01100111/Tickermaster/internal/app/otpChallenge"
	"github.com/Kushian01100111/Tickermaster/internal/app/user"
	"github.com/Kushian01100111/Tickermaster/internal/app/venue"
	"github.com/Kushian01100111/Tickermaster/internal/config"
	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
	http1 "github.com/Kushian01100111/Tickermaster/internal/http"
	"github.com/Kushian01100111/Tickermaster/internal/http/handlers"
	"github.com/Kushian01100111/Tickermaster/internal/http/middleware"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"github.com/Kushian01100111/Tickermaster/internal/storage/mongodb"
	redisDB "github.com/Kushian01100111/Tickermaster/internal/storage/redisdb"
)

/*
	-> Idea sobre compras de tickets
	Estado compartido entre compradores manejado de forma externa a Gin(debido al comportamiento de sync.Pools utilizadas dentro de gin es imposible conocer el estado de una llamada concurrente a otra en http) con holds temporales a traves del cache(redis) en donde todos los compradores comparten e interacturan sobre el mismo estado de la aplicación.
*/

func main() {
	var db *mongo.Client
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	addr := flag.String("addr", ":"+config.Port, "HTTP network address")

	// MongoDB
	db, err = mongodb.ConnectDB(config.DSN, config.DB)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// RedisDB
	rdb, err := redisDB.ConnectRDB(config.RDBSecrets)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		_ = db.Disconnect(context.Background())
		_ = rdb.Close()
	}()

	JWTManager, err := session.NewJWTManager(config.JWTSecrets)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	hasher := session.NewBcryptHasher(0)
	authMiddleware := middleware.NewAuthMiddleware(JWTManager)

	//Entities repositories
	eventRepo := repository.NewEventRepository(db.Database(config.DB))
	venueRepo := repository.NewVenueRepository(db.Database(config.DB))
	userRepo := repository.NewUserRepository(db.Database(config.DB))

	//Non-entities repositories
	authRepo := repository.NewAuthRepository(db.Database(config.DB))
	otpRepo := repository.NewOTPRepository(rdb)
	bookingRepo := repository.NewBookingRepo(rdb)
	emailRepo := email.NewEmailSender(config.ResendAPIKey, config.EmailFrom)

	// Entities service logic
	eventSvc := event.NewEventService(eventRepo, venueRepo)
	venueSvc := venue.NewVenueService(venueRepo)
	userSvc := user.NewUserService(userRepo)

	// Non-entities service logic
	bookingSrv := booking.NewBookingService(bookingRepo)
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

	// handlers
	bookingHandler := handlers.NewBookingHandler(bookingSrv)
	authHandler := handlers.NewAuthHandler(authSrv)
	eventHandler := handlers.NewEventHandler(eventSvc)
	venueHandler := handlers.NewVenueHandler(venueSvc)
	userHandler := handlers.NewUserHandler(userSvc)

	// Main handler of the application
	r := http1.NewHandler(http1.RouterDep{
		AuthHandler:    authHandler,
		BookingHandler: bookingHandler,
		EventDep:       eventHandler,
		VenueDep:       venueHandler,
		UserDep:        userHandler},
		config,
		middleware.Logger(),
		authMiddleware)

	// Server <- is missing some stuff
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
