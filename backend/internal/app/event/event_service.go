package event

import (
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/event"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

type EventParams struct {
	Title       string
	Description string

	Date      *time.Time
	SalesStar *time.Time

	Currency  string
	EventType string
	SeatType  string

	VenueID    uint32
	Performers []string
	Status     string
	Visibility string
}

type SearchParams struct {
	Q        string
	DateForm *time.Time
	DateTo   *time.Time

	Currency string

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

func (s *eventService) SearchEvent(params SearchParams) ([]event.Event, error)
func (s *eventService) GetEvent(name string) (*event.Event, error)
func (s *eventService) CreateEvent(params EventParams) (*event.Event, error)
func (s *eventService) UpdateEvent(name string, params EventParams) (*event.Event, error)
func (s *eventService) DeleteEvent(name string) (*event.Event, error)
