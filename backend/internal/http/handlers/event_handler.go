package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/app/event"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/Kushian01100111/Tickermaster/internal/http/middleware"

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
		context.PUT("", middleware.RequireRole("editor", "admin"), e.createEvent)
		context.GET("", middleware.RequireRole("admin"), e.getAllEvents)
		context.GET("/:id", e.getEvent)
		context.PATCH("/:id", middleware.RequireRole("editor", "admin"), e.updateEvent)
		context.DELETE("/:id", middleware.RequireRole("editor", "admin"), e.deleteEvent)
		context.GET("/search", e.searchEvents)
	}
}

func (e *EventHandler) createEvent(g *gin.Context) {
	var req *dto.EventRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
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

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	event, err := e.app.GetEvent(id, ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToEventResponse(event))
}

func (e *EventHandler) getAllEvents(g *gin.Context) {
	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
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

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
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

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	if err := e.app.DeleteEvent(id, ctx); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.Status(http.StatusNoContent)
}

func (e *EventHandler) searchEvents(g *gin.Context) {
	req, err := ProcessQueries(g)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	events, err := e.app.SearchEvent(event.SearchParams{
		Tokens:       req.Query, // Tener en cuenta en el futuro searchTexts noramlizados para todos los eventos
		DateForm:     req.DateFrom,
		DateTo:       req.DateTo,
		Currency:     req.Currency,
		VenueID:      req.VenueID,
		Availability: req.Availability,
		SortBy:       req.SortBy,
		SortDir:      req.SortDir,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, dto.ToEventResponseSlice(events))
}

func ProcessQueries(g *gin.Context) (*dto.EventSearchRequest, error) {
	var req dto.EventSearchRequest
	layout := "2006-01-02 15:04:05"

	if q := g.Query("q"); q != "" {
		req.Query = getTokens(q)
	}

	if q := g.Query("from"); q != "" {
		date, err := time.Parse(layout, q)
		if err != nil {
			return nil, err
		}
		req.DateFrom = date
	}

	if q := g.Query("to"); q != "" {
		date, err := time.Parse(layout, q)
		if err != nil {
			return nil, err
		}
		req.DateTo = date
	}

	if q := g.Query("currency"); q != "" {
		req.Currency = q
	}

	if q := g.Query("venue"); q != "" {
		req.VenueID = q
	}

	if q := g.Query("availability"); q != "" {
		req.Availability = q
	}

	if q := g.Query("sortBy"); q != "" {
		req.SortBy = q
	}

	if q := g.Query("sortDir"); q != "" {
		number, err := strconv.ParseInt(q, 10, 64)
		if err != nil {
			return nil, err
		}
		req.SortDir = int(number)
	}

	return &req, nil
}

func getTokens(str string) []string {
	res := strings.TrimSpace(str)
	res = strings.ToLower(str)

	res = strings.NewReplacer("-", " ", "_", " ", ".", " ", ",", " ").Replace(res)

	return strings.Fields(res)
}
