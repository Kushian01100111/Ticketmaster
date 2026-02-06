package dto

import (
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
)

//Event

type EventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`

	StartingDate time.Time `json:"startingDate"`
	SalesStart   time.Time `json:"salesStart"`
	Currency     string    `json:"currency"`

	EventType string `json:"eventType"`
	SeatType  string `json:"seatType"`

	VenueID    string   `json:"venue"`
	Performers []string `json:"artist,omitempty"`
	Visibility string   `json:"visibility"`
}

type InternalVenueResponse struct {
	Name     string `json:"name"`
	Address  string `json:"addresss"`
	Capacity int32  `json:"capacity"`
}

type EventResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Date       time.Time `json:"startingDate"`
	SalesStart time.Time `json:"salesStart"`

	Currency  string `json:"currency"`
	EventType string `json:"eventType"`
	SeatType  string `json:"SeatType"`

	VenueID string                `json:"venueId"`
	Venue   InternalVenueResponse `json:"venue"`

	Performers []string `json:"artist,omitempty"`
	Status     string   `json:"status"`
	Visibility string   `json:"visibility,omitempty"`
}

func ToEventResponse(event *event.Event) EventResponse {
	return EventResponse{
		ID:          event.ID.Hex(),
		Title:       event.Title,
		Description: event.Description,
		Date:        event.Date,
		SalesStart:  event.SalesStart,
		EventType:   event.EventType,
		SeatType:    event.SeatType,
		VenueID:     event.VenueID.Hex(),
		Venue:       InternalVenueResponse(event.Venue),
		Performers:  event.Performers,
		Visibility:  event.Visibility,
	}
}

/// Venue

type VenueRequest struct {
	Name string `json:"name"`

	SeatType  string `json:"seatType"`
	SeatMapID string `json:"seatMapId"`

	Address  string `json:"address"`
	Capacity int32  `json:"capacity"`
}

type VenueResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	SeatType  string `json:"seatType"`
	SeatMapID string `json:"seatMapId"`

	Address  string `json:"address"`
	Capacity int32  `json:"capacity"`
}

func ToVenueResponse(venue *venue.Venue) VenueResponse {
	return VenueResponse{
		ID:        venue.ID.Hex(),
		Name:      venue.Name,
		SeatType:  venue.SeatType,
		SeatMapID: venue.SeatMapID.Hex(),
		Address:   venue.Address,
		Capacity:  venue.Capacity,
	}
}
