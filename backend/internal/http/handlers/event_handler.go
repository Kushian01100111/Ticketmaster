package handlers

import (
	"net/http"
	"strings"

	"github.com/Kushian01100111/Tickermaster/internal/app/event"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	app event.EventService
}

func NewEventHandler(svc event.EventService) *EventHandler {
	return &EventHandler{app: svc}
}

func (e *EventHandler) EventRoutes(r *gin.RouterGroup) {
	context := r.Group("/event")
	{
		context.PUT("", e.createEvent)
		context.GET("", e.getAllEvents)
		context.GET("/:id", e.getEvent)
		context.PATCH("/:id", e.updateEvent)
		context.DELETE("/:id", e.deleteEvent)
		context.GET("/search/:id", e.searchEvents)
	}
}

func (e *EventHandler) createEvent(g *gin.Context) {
	var req *dto.EventRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	ctx, cancel := generateCtx()
	defer cancel()

	event, err := e.app.CreateEvent(event.EventParams{
		Title:             req.Title,
		Description:       req.Description,
		StartingDate:      req.StartingDate,
		SalesStartingDate: req.SalesStartingDate,
		Currency:          req.Currency,
		EventType:         req.EventType,
		SeatType:          req.SeatType,
		VenueID:           req.VenueID,
		Performers:        req.Performers,
		Status:            "draft",
		Availability:      req.Availability,
		Visibility:        req.Visibility,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, dto.ToEventResponse(event))
}

func (e *EventHandler) getEvent(g *gin.Context) {
	id := g.Param("id")

	ctx, cancel := generateCtx()
	defer cancel()

	event, err := e.app.GetEvent(id, ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToEventResponse(event))
}

func (e *EventHandler) getAllEvents(g *gin.Context) {
	ctx, cancel := generateCtx()
	defer cancel()

	events, err := e.app.GetAllEvents(ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToEventResponseSlice(events))
}

func (e *EventHandler) updateEvent(g *gin.Context) {
	var req *dto.EventRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	id := g.Param("id")

	ctx, cancel := generateCtx()
	defer cancel()

	event, err := e.app.UpdateEvent(id, event.EventParams{
		Title:             req.Title,
		Description:       req.Description,
		StartingDate:      req.StartingDate,
		SalesStartingDate: req.SalesStartingDate,
		Currency:          req.Currency,
		EventType:         req.EventType,
		SeatType:          req.SeatType,
		VenueID:           req.VenueID,
		Performers:        req.Performers,
		Status:            req.Status,
		Availability:      req.Availability,
		Visibility:        req.Visibility,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, dto.ToEventResponse(event))
}
func (e *EventHandler) deleteEvent(g *gin.Context) {
	id := g.Param("id")

	ctx, cancel := generateCtx()
	defer cancel()

	if err := e.app.DeleteEvent(id, ctx); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.Status(http.StatusNoContent)
}

func (e *EventHandler) searchEvents(g *gin.Context) {
	name := g.Param("name")
	name = DeSlash(name)
}

func DeSlash(str string) string {
	var res strings.Builder
	res.Grow(len(str))

	for _, char := range str {
		if char == '-' {
			res.WriteRune(' ')
		} else {
			res.WriteRune(char)
		}
	}

	return res.String()
}
