package event

import (
	"errors"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

var (
	ErrValidation   = errors.New("string validation error")
	ErrEventType    = errors.New("invalid event type")
	ErrSeatType     = errors.New("invalid seat type")
	ErrAvailability = errors.New("invalid availability type")
)

type EventParams struct {
	Title       string
	Description string

	Date       time.Time
	SalesStart time.Time

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
	SearchEvent(search SearchParams) ([]event.Event, error)
	GetEvent(name string) (*event.Event, error)
	CreateEvent(params EventParams) (*event.Event, error)
	UpdateEvent(name string, params EventParams) (*event.Event, error)
	DeleteEvent(name string) (*event.Event, error)
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{
		eventRepo: repo,
	}
}

func (s *eventService) SearchEvent(params SearchParams) ([]event.Event, error) // Missing: Availability in searchParams
func (s *eventService) GetEvent(name string) (*event.Event, error)

func (s *eventService) CreateEvent(params EventParams) (*event.Event, error) {
	if err := validateParam(params); err != nil {
		return nil, err
	}

}

func (s *eventService) UpdateEvent(name string, params EventParams) (*event.Event, error)
func (s *eventService) DeleteEvent(name string) (*event.Event, error)

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

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusPostpond  EventStatus = "postpond"
)

func validateEventType(str string) error {
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
