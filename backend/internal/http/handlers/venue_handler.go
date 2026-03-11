package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/app/venue"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	app venue.VenueService
}

func NewVenueHandler(svc venue.VenueService) *VenueHandler {
	return &VenueHandler{
		app: svc,
	}
}

func (v *VenueHandler) VenueRoutes(r *gin.RouterGroup) {
	context := r.Group("/venue")
	{
		context.GET("", v.getAllvenues)
		context.GET("/:id", v.getVenue)
		context.PATCH("/:id", v.updateVenue)
		context.DELETE("/:id", v.deleteVenue)
		context.PUT("", v.createVenue)
	}
}

func (v *VenueHandler) getAllvenues(g *gin.Context) {
	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	venues, err := v.app.GetAllVenues(ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToVenueSliceResponse(venues))
}

func (v *VenueHandler) getVenue(g *gin.Context) {
	id := g.Param("id")

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	venue, err := v.app.GetVenue(id, ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToVenueResponse(venue))
}

func (v *VenueHandler) updateVenue(g *gin.Context) {
	var req *dto.VenueRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	id := g.Param("id")

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	venue, err := v.app.UpdateVenue(id, venue.VenueParams{
		Name:      req.Name,
		SeatType:  req.SeatType,
		SeatMapID: req.SeatMapID,
		Address:   req.Address,
		Capacity:  req.Capacity,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, dto.ToVenueResponse(venue))
}

func (v *VenueHandler) deleteVenue(g *gin.Context) {
	id := g.Param("id")

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	if err := v.app.DeleteVenue(id, ctx); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.Status(http.StatusNoContent)
}

func (v *VenueHandler) createVenue(g *gin.Context) {
	var req *dto.VenueRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	venue, err := v.app.CreateVenue(venue.VenueParams{
		Name:      req.Name,
		SeatType:  req.SeatType,
		SeatMapID: req.SeatMapID,
		Address:   req.Address,
		Capacity:  req.Capacity,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, dto.ToVenueResponse(venue))
}

func generateCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}
