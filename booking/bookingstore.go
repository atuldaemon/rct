package booking

import (
	"errors"
	"sync"
	"time"
)

type BookingStore interface {
	Book(spotId int, startTime time.Time, duration time.Duration) (Booking, error)
	Delete(bookingId int) error
	Find(bookingId int) (Booking, error)
	GetAll() ([]Booking, error)
}

type Booking struct {
	ID        int           `json:"id"`
	SpotId    int           `json:"spotId"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrNotFound        = errors.New("not found")
	ErrInvalidReq      = errors.New("invalid request")
	ErrInternal        = errors.New("internal data error")
)

type InMemStore struct {
	mtx   sync.RWMutex
	m     map[int]Booking
	nxtId int // keeps track of the id of the next element to be created
}

func NewInMemBookingStore() (BookingStore, error) {
	s := &InMemStore{m: make(map[int]Booking, 0), nxtId: 1}
	return s, nil
}

func (s *InMemStore) Book(spotId int, startTime time.Time, duration time.Duration) (Booking, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	b := Booking{ID: s.nxtId, SpotId: spotId, StartTime: startTime, Duration: duration}
	s.m[b.ID] = b
	s.nxtId++
	return b, nil
}

func (s *InMemStore) Delete(bookingId int) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.m, bookingId)
	return nil
}

func (s *InMemStore) Find(bookingId int) (Booking, error) {
	b, ok := s.m[bookingId]
	if !ok {
		return Booking{}, ErrNotFound
	}
	return b, nil
}

func (s *InMemStore) GetAll() ([]Booking, error) {
	bb := make([]Booking, 0)
	for _, b := range s.m {
		bb = append(bb, b)
	}
	return bb, nil
}
