package handlers

import (
	"net/http"

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
		context.GET("/:name", v.getVenue)
		context.PUT("", v.createVenue)
	}
}

func (v *VenueHandler) getVenue(g *gin.Context) {

}

func (v *VenueHandler) createVenue(g *gin.Context) {
	var req *dto.VenueRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	venue, err := v.app.CreateVenue(venue.VenueParams{
		Name:      req.Name,
		SeatType:  req.Address,
		SeatMapID: req.SeatMapID,
		Address:   req.Address,
		Capacity:  req.Capacity,
	})

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, dto.ToVenueResponse(venue))
}
