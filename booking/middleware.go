package booking

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

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

func (mw loggingMiddleware) GetAll(ctx context.Context) (b []Booking, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAll", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAll(ctx)
}

func (mw loggingMiddleware) Book(ctx context.Context, spotId string, startTime time.Time, duration time.Duration) (b Booking, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Book", "spotId", spotId, "startTime", startTime, "duration", duration, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Book(ctx, spotId, startTime, duration)
}

func (mw loggingMiddleware) Delete(ctx context.Context, bookingId string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Delete", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Delete(ctx, bookingId)
}
