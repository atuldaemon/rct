package parking

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"context"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) GetAll(ctx context.Context) ([]Spot, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetAll").Add(1)
		s.requestLatency.With("method", "GetAll").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAll(ctx)
}


func (s *instrumentingService) GetFree(ctx context.Context) ([]Spot, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetFree").Add(1)
		s.requestLatency.With("method", "GetFree").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetFree(ctx)
}


func (s *instrumentingService) GetReserved(ctx context.Context) ([]Spot, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetReserved").Add(1)
		s.requestLatency.With("method", "GetReserved").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetReserved(ctx)
}

func (s *instrumentingService) Search(ctx context.Context, lat, lon, radius string, metric SearchMetric) ([]ExtendedSpot, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Search").Add(1)
		s.requestLatency.With("method", "Search").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Search(ctx, lat, lon, radius, metric)
}

func (s *instrumentingService) FindById(ctx context.Context, id string) (Spot, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "FindById").Add(1)
		s.requestLatency.With("method", "FindById").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.FindById(ctx, id)
}

func (s *instrumentingService) Update(ctx context.Context, sp Spot) (Spot, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Update").Add(1)
		s.requestLatency.With("method", "Update").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Update(ctx, sp)
}
