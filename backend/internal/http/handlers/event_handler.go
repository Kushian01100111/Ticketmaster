package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

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
		context.GET("", e.searchEvents)
		context.GET("/:name", e.getEvent)
		context.PATCH("/:name", e.updateEvent)
		context.PUT("", e.createEvent)
		context.DELETE("/:mane", e.deleteEvent)
	}
}

func (e *EventHandler) createEvent(g *gin.Context) {
	var req *dto.EventRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	event, err := e.app.CreateEvent(event.EventParams{
		Title:       req.Title,
		Description: req.Description,
		Date:        req.StartingDate,
		SalesStart:  req.SalesStart,
		Currency:    req.Currency,
		EventType:   req.EventType,
		SeatType:    req.SeatType,
		VenueID:     req.VenueID,
		Performers:  req.Performers,
		Status:      "draft",
		Visibility:  req.Visibility,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	g.JSON(http.StatusCreated, dto.ToEventResponse(event))
}

func (e *EventHandler) searchEvents(g *gin.Context) {
	name := g.Param("name")
	name = deSlash(name)
}
func (e *EventHandler) getEvent(g *gin.Context) {
	name := g.Param("name")
	name = deSlash(name)

}

func (e *EventHandler) updateEvent(g *gin.Context) {
	name := g.Param("name")
	name = deSlash(name)
}
func (e *EventHandler) deleteEvent(g *gin.Context) {
	name := g.Param("name")
	name = deSlash(name)
}

func deSlash(str string) string {
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
