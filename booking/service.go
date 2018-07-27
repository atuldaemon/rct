package booking

import (
	"context"
	"errors"
	"time"

	"strconv"

	"github.com/atuldaemon/rct/parking"
)

var (
	ErrInvalidSpotId             = errors.New("invalid spotid passed in booking request")
	ErrAlreadyReserved           = errors.New("spot already reserved")
	ErrInvalidBookingId          = errors.New("invalid booking id  passed in delete booking request")
	ErrInvalidSpotIdForBookingId = errors.New("invalid slot id for booking id passed in delete booking request")
	ErrFailedToUpdate            = errors.New("Failed to update/release slot")
)

type Service interface {
	GetAll(ctx context.Context) ([]Booking, error)
	Book(ctx context.Context, spotId string, startTime time.Time, duration time.Duration) (Booking, error)
	Delete(ctx context.Context, bookingId string) error
}

type service struct {
	bookingStore   BookingStore
	parkingService parking.Service
}

func NewService(bookingStore BookingStore, pService parking.Service) Service {
	return &service{bookingStore: bookingStore, parkingService: pService}
}

func (s *service) GetAll(ctx context.Context) ([]Booking, error) {
	return s.bookingStore.GetAll()
}

func (s *service) Book(ctx context.Context, spotId string, startTime time.Time, duration time.Duration) (Booking, error) {
	spot, err := s.parkingService.FindById(ctx, string(spotId))
	if err != nil {
		return Booking{}, ErrInvalidSpotId
	}
	if spot.IsReserved == true {
		return Booking{}, ErrAlreadyReserved
	}
	spotIdInt, err := strconv.Atoi(spotId)
	if err != nil {
		return Booking{}, ErrInvalidReq
	}
	spot.IsReserved = true
	_, err = s.parkingService.Update(ctx, spot)
	if err != nil {
		return Booking{}, ErrInternal
	}
	return s.bookingStore.Book(spotIdInt, startTime, duration)
}

func (s *service) Delete(ctx context.Context, bookingId string) error {
	bookingIdInt, err := strconv.Atoi(bookingId)
	if err != nil {
		return ErrInvalidReq
	}
	b, err := s.bookingStore.Find(bookingIdInt)
	if err != nil {
		return ErrInvalidBookingId
	}
	spot, err := s.parkingService.FindById(ctx, strconv.Itoa(b.SpotId))
	if err != nil {
		return ErrInvalidSpotIdForBookingId
	}
	spot.IsReserved = false
	_, err = s.parkingService.Update(ctx, spot)
	if err != nil {
		return ErrFailedToUpdate
	}
	return s.bookingStore.Delete(bookingIdInt)
}
