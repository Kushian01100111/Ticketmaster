package handlers

import (
	"github.com/Kushian01100111/Tickermaster/internal/app/booking"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	app booking.BookingService
}

func NewBookingHandler(svc booking.BookingService) *BookingHandler {
	return &BookingHandler{app: svc}
}

func (e *BookingHandler) PublicRoutes(r *gin.RouterGroup) {
}

func (e *BookingHandler) PrivateRoutes(r *gin.RouterGroup) {

}
