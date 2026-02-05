package venue

import (
	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

type VenueParams struct {
	Name string

	SeatType  string
	SeatMapID string

	Address  string
	Capacity int32
}

type VenueService interface {
	GetVenue(venueID string) (*venue.Venue, error)
	CreateVenue(params VenueParams) (*venue.Venue, error)
	UpdateVenue(venueID string, params VenueParams) (*venue.Venue, error)
	DeleteVenue(venueID string) error
}

type venueService struct {
	venueRepo repository.VenueRepository
}

func NewVenueService(repo repository.VenueRepository) VenueService {
	return &venueService{
		venueRepo: repo,
	}
}

func (s *venueService) GetVenue(venueID string) (*venue.Venue, error)
func (s *venueService) CreateVenue(params VenueParams) (*venue.Venue, error)
func (s *venueService) UpdateVenue(venueID string, params VenueParams) (*venue.Venue, error)
func (s *venueService) DeleteVenue(venueID string) error
