package event

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrValidation              = errors.New("string validation error")
	ErrStatus                  = errors.New("invalid event status type")
	ErrEventType               = errors.New("invalid event type")
	ErrSeatType                = errors.New("invalid seat type")
	ErrVisibility              = errors.New("invalid visibility type")
	ErrAvailability            = errors.New("invalid availability type")
	ErrSalesDateWithStartEvent = errors.New("sales date most be at lest a hour prior to the date of the start of the event")
)

type EventParams struct {
	Title       string
	Description string

	StartingDate time.Time
	SalesStart   time.Time

	Currency  string
	EventType string
	SeatType  string

	VenueID      string
	Performers   []string
	Status       string
	Visibility   string
	Availability string
}

type SearchParams struct {
	Q        string
	DateForm time.Time
	DateTo   time.Time

	Currency string

	VenueID      string
	Availability []string
	SortBy       string
	SortDir      int
}

type EventService interface {
	GetEvent(idHex string, ctx context.Context) (*event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
	CreateEvent(params EventParams, ctx context.Context) (*event.Event, error)
	UpdateEvent(eventID string, params EventParams, ctx context.Context) (*event.Event, error)
	DeleteEvent(eventID string, ctx context.Context) error
	SearchEvent(search SearchParams, ctx context.Context) ([]event.Event, error)
}

type eventService struct {
	eventRepo repository.EventRepository
	venueRepo repository.VenueRepository
}

func NewEventService(eventrepo repository.EventRepository, venuerepo repository.VenueRepository) EventService {
	return &eventService{
		eventRepo: eventrepo,
		venueRepo: venuerepo,
	}
}

func (s *eventService) GetEvent(idHex string, ctx context.Context) (*event.Event, error) {
	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	return s.eventRepo.GetByID(id, ctx)
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	return s.eventRepo.GetAllEvents(ctx)
} // <- Estaria faltando un filtro, así asegurar me que solo cierto eventos sean devueltos

func (s *eventService) CreateEvent(params EventParams, ctx context.Context) (*event.Event, error) {
	if err := validateParam(params); err != nil {
		return nil, err
	}

	if err := validateDatesEvent(params.StartingDate, params.SalesStart); err != nil {
		return nil, err
	}

	venueID, err := bson.ObjectIDFromHex(params.VenueID)
	if err != nil {
		return nil, err
	}

	venue, err := s.venueRepo.GetByID(venueID, ctx)
	if err != nil {
		return nil, err
	}

	Event := &event.Event{
		Title:       params.Title,
		Description: params.Description,
		Date:        params.StartingDate,
		SalesStart:  params.SalesStart,
		Currency:    params.Currency,
		EventType:   params.EventType,
		SeatType:    params.SeatType,
		VenueID:     venue.ID,
		Venue: event.Venue{
			Name:     venue.Name,
			Address:  venue.Address,
			Capacity: venue.Capacity,
		},
		Performers:   params.Performers,
		Status:       params.Status,
		Availability: params.Availability,
		Visibility:   params.Visibility,
	}

	if err := s.eventRepo.Create(Event, ctx); err != nil {
		return nil, err
	}

	return Event, nil
}

func (s *eventService) UpdateEvent(name string, params EventParams, ctx context.Context) (*event.Event, error) {
	return nil, nil
}

func (s *eventService) DeleteEvent(name string, ctx context.Context) error {
	return nil
}

func (s *eventService) SearchEvent(params SearchParams, ctx context.Context) ([]event.Event, error) {
	var res []event.Event
	return res, nil
}

///
///
///
///

func validateDatesEvent(startEvent time.Time, startSales time.Time) error {
	if startSales.Before(startEvent) || startEvent.Sub(startSales) >= time.Hour {
		return ErrSalesDateWithStartEvent
	}
	return nil
}

func validateParam(params EventParams) error {
	if err := validateString(params.Title); err != nil {
		return err
	}

	if err := validateString(params.Description); err != nil {
		return err
	}

	if err := validateString(params.Currency); err != nil {
		return err
	}

	if err := validateEventType(params.EventType); err != nil {
		return err
	}

	if err := validateSeatType(params.SeatType); err != nil {
		return err
	}

	if err := validateEventStatus(params.Status); err != nil {
		return err
	}

	if err := validateVisibility(params.Visibility); err != nil {
		return err
	}

	if err := validateAvailabity(params.Availability); err != nil {
		return err
	}

	return nil
}

func validateString(str string) error {
	if strings.TrimSpace(str) == "" {
		return ErrValidation
	}
	return nil
}

type EventType string

const (
	EventConsert EventType = "concert"
	EventRecital EventType = "recital"
	EventSolo    EventType = "solo recital"
	EventOpera   EventType = "operatic productions"
)

func validateEventType(str string) error {
	switch EventType(str) {
	case EventConsert, EventRecital, EventSolo, EventOpera:
		return nil
	default:
		return ErrEventType
	}
}

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusPostpond  EventStatus = "postpond"
)

func validateEventStatus(str string) error {
	switch EventStatus(str) {
	case EventStatusDraft, EventStatusPublished, EventStatusCancelled, EventStatusPostpond:
		return nil
	default:
		return ErrEventType
	}
}

type SeatType string

const (
	SeatTypeSeated   SeatType = "seated"
	SeatTypeStanding SeatType = "standing"
	SeatTypeMixed    SeatType = "mixed"
)

func validateSeatType(str string) error {
	switch SeatType(str) {
	case SeatTypeSeated, SeatTypeStanding, SeatTypeMixed:
		return nil
	default:
		return ErrSeatType
	}
}

type Visibility string

const (
	VisibilityPublic   Visibility = "public"
	VisibilityUnlisted Visibility = "unlisted"
	VisibilityPrivate  Visibility = "private"
)

func validateVisibility(str string) error {
	switch Visibility(str) {
	case VisibilityPublic, VisibilityUnlisted, VisibilityPrivate:
		return nil
	default:
		return ErrVisibility
	}
}

type Availability string

const (
	AvailabilityAvailable Availability = "available"
	AvailabilitySoldOut   Availability = "soldOut"
)

func validateAvailabity(str string) error {
	switch Availability(str) {
	case AvailabilityAvailable, AvailabilitySoldOut:
		return nil
	default:
		return ErrAvailability
	}
}
