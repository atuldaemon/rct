package booking

import (
	"time"

	"context"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) GetAll(ctx context.Context) (b []Booking, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetAll").Add(1)
		s.requestLatency.With("method", "GetAll").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAll(ctx)
}

func (s *instrumentingService) Book(ctx context.Context, spotId string, startTime time.Time, duration time.Duration) (Booking, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Book").Add(1)
		s.requestLatency.With("method", "Book").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Book(ctx, spotId, startTime, duration)
}

func (s *instrumentingService) Delete(ctx context.Context, bookingId string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Delete").Add(1)
		s.requestLatency.With("method", "Delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Delete(ctx, bookingId)
}
