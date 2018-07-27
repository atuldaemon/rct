package parking

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetAll(ctx context.Context) (sp []Spot, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAllParking", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAll(ctx)
}

func (mw loggingMiddleware) GetFree(ctx context.Context) (sp []Spot, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetFreeParking", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetFree(ctx)
}

func (mw loggingMiddleware) GetReserved(ctx context.Context) (sp []Spot, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetReservedParking", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetReserved(ctx)
}

func (mw loggingMiddleware) Search(ctx context.Context, lat, lon, rad string, metric SearchMetric) (sp []ExtendedSpot, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Search", "lat", lat, "lon", lon, "radius", rad, "metric", metric, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Search(ctx, lat, lon, rad, metric)
}

func (mw loggingMiddleware) FindById(ctx context.Context, id string) (sp Spot, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Find", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.FindById(ctx, id)
}

func (mw loggingMiddleware) Update(ctx context.Context, s Spot) (sp Spot, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Update", "id", s.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Update(ctx, s)
}
