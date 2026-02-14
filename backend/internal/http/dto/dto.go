package dto

import (
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
)

//Event

type EventSearchRequest struct {
	Name     string    `json:"name"`
	DateFrom time.Time `json:"dateFrom,omitempty"`
	DateTo   time.Time `json:"dateTo,omitempty"`

	Currency string `json:"currency,omitempty"`

	VenueID      string `json:"venueId,omitempty"`
	Availability string `json:"availability,omitempty"`
	SortBy       string `json:"sortBy,omitempty"`
	SortDir      int    `json:"sortDir,omitempty"`
}

type EventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`

	StartingDate      time.Time `json:"startingDate"`
	SalesStartingDate time.Time `json:"salesStartingDate"`
	Currency          string    `json:"currency"`

	EventType string `json:"eventType"`
	SeatType  string `json:"seatType"`

	VenueID      string   `json:"venueId"`
	Performers   []string `json:"performers,omitempty"`
	Status       string   `json:"status,omitempty"`
	Availability string   `json:"availability"`
	Visibility   string   `json:"visibility"`
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

	StartingDate      time.Time `json:"startingDate"`
	SalesStartingDate time.Time `json:"salesStartingDate"`

	Currency  string `json:"currency"`
	EventType string `json:"eventType"`
	SeatType  string `json:"seatType"`

	VenueID string                `json:"venueId"`
	Venue   InternalVenueResponse `json:"venue"`

	Performers   []string `json:"performers,omitempty"`
	Status       string   `json:"status"`
	Availability string   `json:"availability"`
	Visibility   string   `json:"visibility,omitempty"`
}

func ToEventResponse(event *event.Event) EventResponse {
	return EventResponse{
		ID:                event.ID.Hex(),
		Title:             event.Title,
		Description:       event.Description,
		StartingDate:      event.StartingDate,
		SalesStartingDate: event.SalesStartingDate,
		Currency:          event.Currency,
		EventType:         event.EventType,
		SeatType:          event.SeatType,
		VenueID:           event.VenueID.Hex(),
		Venue:             InternalVenueResponse(event.Venue),
		Performers:        event.Performers,
		Status:            event.Status,
		Availability:      event.Availability,
		Visibility:        event.Visibility,
	}
}

func ToEventResponseSlice(events []event.Event) []EventResponse {
	response := make([]EventResponse, len(events))

	for i, event := range events {
		response[i] = ToEventResponse(&event)
	}
	return response
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

func ToVenueSliceResponse(venues []venue.Venue) []VenueResponse {
	response := make([]VenueResponse, len(venues))

	for i, venue := range venues {
		response[i] = ToVenueResponse(&venue)
	}

	return response
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
