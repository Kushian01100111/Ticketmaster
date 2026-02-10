package venue

import (
	"context"
	"errors"
	"strings"

	"github.com/Kushian01100111/Tickermaster/internal/domain/venue"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrValidation    = errors.New("string validation error")
	ErrSeatType      = errors.New("invalid seat type")
	ErrCapacity      = errors.New("invalid capacity")
	ErrProvidedID    = errors.New("provided id is not a valid objectID")
	ErrProvidedMapID = errors.New("provided mapId is no a valid objectId")
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
	GetAllVenues(ctx context.Context) ([]venue.Venue, error)
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
	id, err := bson.ObjectIDFromHex(venueID)
	if err != nil {
		return nil, err
	}
	return s.venueRepo.GetByID(id, ctx)
}

func (s *venueService) CreateVenue(params VenueParams, ctx context.Context) (*venue.Venue, error) {
	if err := validateParam(params); err != nil {
		return nil, err
	}

	SeatMapID := bson.NewObjectID()

	/*
	   	-> Hasta que no pueda asegurar que los ID seatMap sean correctos lo cambio a una generado de forma random
	   SeatMapID, err := primitive.ObjectIDFromHex(params.SeatMapID)
	   	if err != nil {
	   		return nil, err
	   	}
	*/

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

func (s *venueService) GetAllVenues(ctx context.Context) ([]venue.Venue, error) {
	return s.venueRepo.GetAll(ctx)
}

func (s *venueService) UpdateVenue(venueID string, params VenueParams, ctx context.Context) (*venue.Venue, error) {
	if err := validateParam(params); err != nil {
		return nil, ErrValidation
	}

	id, err := bson.ObjectIDFromHex(venueID)
	if err != nil {
		return nil, ErrProvidedID
	}

	Venue := &venue.Venue{
		ID:       id,
		Name:     params.Name,
		SeatType: params.SeatType,
		Address:  params.Address,
		Capacity: params.Capacity,
	}

	if params.SeatMapID != "" {
		SeatMapID, err := bson.ObjectIDFromHex(params.SeatMapID)
		if err != nil {
			return nil, ErrProvidedMapID
		}
		Venue.SeatMapID = SeatMapID
	}

	if err := s.venueRepo.Update(Venue, ctx); err != nil {
		return nil, err
	}

	return Venue, nil
}

func (s *venueService) DeleteVenue(venueID string, ctx context.Context) error {
	id, err := bson.ObjectIDFromHex(venueID)
	if err != nil {
		return err
	}

	if err := s.venueRepo.Delete(id, ctx); err != nil {
		return err
	}

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
