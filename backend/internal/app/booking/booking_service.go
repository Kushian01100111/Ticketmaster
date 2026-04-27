package booking

import "github.com/Kushian01100111/Tickermaster/internal/repository"

type BookingService interface {
}

type bookingService struct {
	bookingRepo repository.BookingRepo
}

func NewBookingService(repo repository.BookingRepo) BookingService {
	return &bookingService{bookingRepo: repo}
}
