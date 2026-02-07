package venue

import (
	"context"
	"errors"
	"strings"

	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrValidation = errors.New("string validation error")
	ErrSeatType   = errors.New("invalid seat type")
	ErrCapacity   = errors.New("invalid capacity")
)

type VenueParams struct {
	Name string

	SeatType  string
	SeatMapID string

	Address  string
	Capacity int32
}

type VenueService interface {
	GetVenue(venueID string, ctx context.Context) (*venue.Venue, error)
	CreateVenue(params VenueParams, ctx context.Context) (*venue.Venue, error)
	UpdateVenue(venueID string, params VenueParams, ctx context.Context) (*venue.Venue, error)
	DeleteVenue(venueID string, ctx context.Context) error
}

type venueService struct {
	venueRepo repository.VenueRepository
}

func NewVenueService(repo repository.VenueRepository) VenueService {
	return &venueService{
		venueRepo: repo,
	}
}

func (s *venueService) GetVenue(venueID string, ctx context.Context) (*venue.Venue, error) {
	return s.venueRepo.GetByID(venueID, ctx)
}

func (s *venueService) CreateVenue(params VenueParams, ctx context.Context) (*venue.Venue, error) {
	if err := validateParam(params); err != nil {
		return nil, err
	}

	SeatMapID, err := primitive.ObjectIDFromHex(params.SeatMapID)
	if err != nil {
		return nil, err
	}

	Venue := &venue.Venue{
		Name:      params.Name,
		SeatType:  params.SeatType,
		SeatMapID: SeatMapID,
		Address:   params.Address,
		Capacity:  params.Capacity,
	}

	id, err := s.venueRepo.Create(Venue, ctx)
	if err != nil {
		return nil, err
	}

	Venue.ID = id

	return Venue, nil
}

func (s *venueService) UpdateVenue(venueID string, params VenueParams, ctx context.Context) (*venue.Venue, error) {
	return nil, nil
}

func (s *venueService) DeleteVenue(venueID string, ctx context.Context) error {
	return nil
}

func validateParam(params VenueParams) error {
	if err := validateString(params.Name); err != nil {
		return err
	}

	if err := validateSeatType(params.SeatType); err != nil {
		return err
	}

	if err := validateString(params.Address); err != nil {
		return err
	}

	if params.Capacity < 0 || params.Capacity > 90000 {
		return ErrCapacity
	}

	return nil
}

func validateString(str string) error {
	if strings.TrimSpace(str) == "" {
		return ErrValidation
	}
	return nil
}

type SeatType string

const (
	SeatedType   SeatType = "seated"
	StandingType SeatType = "standing"
	MixedType    SeatType = "mixed"
)

func validateSeatType(str string) error {
	switch SeatType(str) {
	case SeatedType, StandingType, MixedType:
		return nil
	default:
		return ErrSeatType
	}
}
