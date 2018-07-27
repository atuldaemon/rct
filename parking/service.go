package parking

import (
	"context"
	"strconv"
)

// Parking service

type Service interface {
	GetAll(ctx context.Context) ([]Spot, error)
	GetFree(ctx context.Context) ([]Spot, error)
	GetReserved(ctx context.Context) ([]Spot, error)
	Search(ctx context.Context, lat, lon, radius string, metric SearchMetric) ([]ExtendedSpot, error)
	FindById(ctx context.Context, id string) (Spot, error)
	Update(ctx context.Context, sp Spot) (Spot, error)
}

type service struct {
	parkingStore ParkingStore
}

func NewService(store ParkingStore) Service {
	return &service{parkingStore: store}
}

func (s *service) GetAll(ctx context.Context) ([]Spot, error) {
	return s.parkingStore.Get(all)
}

func (s *service) GetFree(ctx context.Context) ([]Spot, error) {
	return s.parkingStore.Get(free)
}

func (s *service) GetReserved(ctx context.Context) ([]Spot, error) {
	return s.parkingStore.Get(reserved)
}

func (s *service) Search(ctx context.Context, lat, lon, radius string, metric SearchMetric) ([]ExtendedSpot, error) {
	return s.parkingStore.Search(lat, lon, radius, metric)
}

func (s *service) FindById(ctx context.Context, id string) (Spot, error) {
	intId, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		return Spot{}, ErrInvalidReq
	}
	return s.parkingStore.FindById(int(intId))
}

func (s *service) Update(ctx context.Context, sp Spot) (Spot, error) {
	return s.parkingStore.Update(sp)
}
