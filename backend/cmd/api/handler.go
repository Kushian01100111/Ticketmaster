package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHandler(config *Config) http.Handler {
	if config.GinConfig == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-type", "Accept", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-length"},
		AllowCredentials: true,
		MaxAge:           25 * time.Minute,
	}))

	return r
}
