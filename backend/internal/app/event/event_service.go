package event

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/text/currency"
)

var (
	ErrValidation              = errors.New("string validation error")
	ErrStatus                  = errors.New("invalid event status type")
	ErrEventType               = errors.New("invalid event type")
	ErrSeatType                = errors.New("invalid seat type")
	ErrVisibility              = errors.New("invalid visibility type")
	ErrAvailability            = errors.New("invalid availability type")
	ErrSalesDateWithStartEvent = errors.New("sale date must be at least one hour before the event start date.")
	ErrProvidedID              = errors.New("provided id is not a valid objectID")
	ErrCurrency                = errors.New("invalid currency type")
	ErrSortBy                  = errors.New("invalid sorted type")
	ErrSortDir                 = errors.New("invalid sortdir type")
)

type EventParams struct {
	Title       string
	Description string

	StartingDate      time.Time
	SalesStartingDate time.Time

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
	Availability string
	SortBy       string
	SortDir      int
}

type EventService interface {
	GetEvent(idHex string, ctx context.Context) (*event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
	CreateEvent(params EventParams, ctx context.Context) (*event.Event, error)
	UpdateEvent(idHex string, params EventParams, ctx context.Context) (*event.Event, error)
	DeleteEvent(idHex string, ctx context.Context) error
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

	venueID, err := bson.ObjectIDFromHex(params.VenueID)
	if err != nil {
		return nil, err
	}

	venue, err := s.venueRepo.GetByID(venueID, ctx)
	if err != nil {
		return nil, err
	}

	Event := &event.Event{
		Title:             params.Title,
		Description:       params.Description,
		StartingDate:      params.StartingDate,
		SalesStartingDate: params.SalesStartingDate,
		Currency:          params.Currency,
		EventType:         params.EventType,
		SeatType:          params.SeatType,
		VenueID:           venue.ID,
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

	id, err := s.eventRepo.Create(Event, ctx)
	if err != nil {
		return nil, err
	}

	Event.ID = id
	return Event, nil
}

func (s *eventService) UpdateEvent(idHex string, params EventParams, ctx context.Context) (*event.Event, error) {
	//fmt.Printf("EventType raw: %q\n", params.EventType)
	if err := validateParam(params); err != nil {
		return nil, err
	}

	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, ErrProvidedID
	}

	venueID, err := bson.ObjectIDFromHex(params.VenueID)
	if err != nil {
		return nil, ErrProvidedID
	}

	venue, err := s.venueRepo.GetByID(venueID, ctx)
	if err != nil {
		return nil, err
	}

	Event := &event.Event{
		ID:                id,
		Title:             params.Title,
		Description:       params.Description,
		StartingDate:      params.StartingDate,
		SalesStartingDate: params.SalesStartingDate,
		Currency:          params.Currency,
		EventType:         params.EventType,
		SeatType:          params.SeatType,
		VenueID:           venue.ID,
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

	if err := s.eventRepo.Update(Event, ctx); err != nil {
		return nil, err
	}

	return Event, nil
}

func (s *eventService) DeleteEvent(idHex string, ctx context.Context) error {
	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	if err := s.eventRepo.Delete(id, ctx); err != nil {
		return err
	}

	return nil
}

func (s *eventService) SearchEvent(params SearchParams, ctx context.Context) ([]event.Event, error) {
	var res []event.Event

	if err := validateSearchParams(params); err != nil {
		return res, nil
	}

	Querie := strings.TrimSpace(params.Q)

	venueID, err := bson.ObjectIDFromHex(params.VenueID)
	if err != nil {
		return res, err
	}

	events, err := s.eventRepo.SearchByParams(&event.SearchEvent{
		Q:            Querie,
		DateFrom:     params.DateForm,
		DateTo:       params.DateTo,
		Currency:     params.Currency,
		VenueID:      venueID,
		Availability: params.Availability,
		SortBy:       params.SortBy,
		SortDir:      params.SortDir,
	}, ctx)
	if err != nil {
		return res, err
	}

	return events, nil
}

///
///
///
///

func validateDatesEvent(startEvent time.Time, startSales time.Time) error {
	if startSales.After(startEvent) || startEvent.Sub(startSales) < time.Hour {
		return ErrSalesDateWithStartEvent
	}
	return nil
}

func validateSearchParams(params SearchParams) error {
	if err := validateString(params.Q); err != nil {
		return err
	}

	if err := validateDatesEvent(params.DateForm, params.DateTo); err != nil {
		return err
	}

	if err := validateCurrency(params.Currency); err != nil {
		return err
	}

	if err := validateString(params.VenueID); err != nil {
		return err
	}

	if err := validateAvailabity(params.Availability); err != nil {
		return err
	}

	if err := validateSortBy(params.SortBy); err != nil {
		return err
	}

	if err := validateSortDir(params.SortDir); err != nil {
		return err
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

	if err := validateDatesEvent(params.StartingDate, params.SalesStartingDate); err != nil {
		return err
	}

	if err := validateCurrency(params.Currency); err != nil {
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
	//fmt.Printf("EventType debug: %q len=%d bytes=%v\n", str, len(str), []byte(str))
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
		return ErrStatus
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

//

func validateCurrency(curr string) error {
	c := strings.ToUpper(strings.TrimSpace(curr))
	if len(c) != 3 {
		return ErrCurrency
	}
	if _, err := currency.ParseISO(c); err != nil {
		return ErrCurrency
	}
	return nil
}

type SortedBy string

const (
	SortedByTitle        SortedBy = "title"
	SortedByDate         SortedBy = "date"
	SortedByAvailability SortedBy = "availability"
	SortedByVenue        SortedBy = "venue"
)

func validateSortBy(sort string) error {
	switch SortedBy(sort) {
	case SortedByTitle, SortedByDate, SortedByAvailability, SortedByVenue:
		return nil
	default:
		return ErrSortBy
	}
}

type SorDir int

const (
	SorDirDesc SorDir = 0
	SorDirAsc  SorDir = 1
)

func validateSortDir(dir int) error {
	switch SorDir(dir) {
	case SorDirAsc, SorDirDesc:
		return nil
	default:
		return ErrSortDir
	}
}
