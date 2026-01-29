package handlers

import (
	"strings"

	"github.com/Kushian01100111/Tickermaster/internal/app/event"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	app event.EventService
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
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

func (e *EventHandler) createEvent(g *gin.Context)
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
