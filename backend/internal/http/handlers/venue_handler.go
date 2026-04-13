package handlers

import (
	"net/http"

	"github.com/Kushian01100111/Tickermaster/internal/app/venue"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/Kushian01100111/Tickermaster/internal/http/middleware"
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

func (v *VenueHandler) PublicRoutes(r *gin.RouterGroup) {
}

func (v *VenueHandler) PrivateRoutes(r *gin.RouterGroup) {
	context := r.Group("/venue")
	{
		context.GET("", middleware.RequireRole("admin"), v.getAllvenues)
		context.GET("/:id", middleware.RequireRole("editor", "admin"), v.getVenue)
		context.PATCH("/:id", middleware.RequireRole("admin"), v.updateVenue)
		context.DELETE("/:id", middleware.RequireRole("admin"), v.deleteVenue)
		context.PUT("", middleware.RequireRole("admin"), v.createVenue)
	}
}

func (v *VenueHandler) getAllvenues(g *gin.Context) {
	venues, err := v.app.GetAllVenues(g.Request.Context())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToVenueSliceResponse(venues))
}

func (v *VenueHandler) getVenue(g *gin.Context) {
	id := g.Param("id")

	venue, err := v.app.GetVenue(id, g.Request.Context())
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

	venue, err := v.app.UpdateVenue(id, venue.VenueParams{
		Name:      req.Name,
		SeatType:  req.SeatType,
		SeatMapID: req.SeatMapID,
		Address:   req.Address,
		Capacity:  req.Capacity,
	}, g.Request.Context())

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, dto.ToVenueResponse(venue))
}

func (v *VenueHandler) deleteVenue(g *gin.Context) {
	id := g.Param("id")

	if err := v.app.DeleteVenue(id, g.Request.Context()); err != nil {
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

	venue, err := v.app.CreateVenue(venue.VenueParams{
		Name:      req.Name,
		SeatType:  req.SeatType,
		SeatMapID: req.SeatMapID,
		Address:   req.Address,
		Capacity:  req.Capacity,
	}, g.Request.Context())

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, dto.ToVenueResponse(venue))
}
