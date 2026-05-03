package http1

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Kushian01100111/Tickermaster/internal/config"
	"github.com/Kushian01100111/Tickermaster/internal/http/handlers"
	"github.com/Kushian01100111/Tickermaster/internal/http/middleware"
)

type RouterDep struct {
	AuthHandler    *handlers.AuthHandler
	BookingHandler *handlers.BookingHandler
	EventDep       *handlers.EventHandler
	VenueDep       *handlers.VenueHandler
	UserDep        *handlers.UserHandler
}

func NewHandler(dep RouterDep, config *config.Config, logger gin.HandlerFunc, auth *middleware.AuthMiddleware) http.Handler {
	if config.GinConfig == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger)
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-type", "Accept", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-length"},
		AllowCredentials: true,
		MaxAge:           25 * time.Minute,
	}))

	api := r.Group("/api")

	//Public
	public := api.Group("")
	{
		dep.AuthHandler.AuthRoutes(public)
		dep.BookingHandler.PublicRoutes(public)
		dep.EventDep.PublicRoutes(public)
		dep.VenueDep.PublicRoutes(public)
		dep.UserDep.PublicRoutes(public)
	}

	//Private
	private := api.Group("")
	private.Use(auth.RequireAuth())
	{
		dep.BookingHandler.PrivateRoutes(private)
		dep.EventDep.PrivateRoutes(private)
		dep.VenueDep.PrivateRoutes(private)
		dep.UserDep.PrivateRoutes(private)
	}

	return r
}
