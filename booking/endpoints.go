package booking

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetAllEndpoint  endpoint.Endpoint
	BookingEndpoint endpoint.Endpoint
	DeleteEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetAllEndpoint:  MakeGetAllEndpoint(s),
		BookingEndpoint: MakeBookingEndpoint(s),
		DeleteEndpoint:  MakeDeleteEndpoint(s),
	}
}

func MakeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bb, e := s.GetAll(ctx)
		return getAllResponse{Bookings: bb, Err: e}, e
	}
}

func MakeBookingEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bookingRequest)
		// TODO: using a default timeslot of 30 mins. Need to take a param
		b, e := s.Book(ctx, req.SpotId, time.Now(), time.Duration(30*time.Minute))
		return bookingResponse{Booking: b, Err: e}, e
	}
}

func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteRequest)
		e := s.Delete(ctx, req.BookingId)
		return deleteResponse{Err: e}, e
	}
}

//

type getAllRequest struct {
}

type bookingRequest struct {
	SpotId string `json:"id"`
	//StartTime time.Time     `json:"startTime"`
	//Duration  time.Duration `json:"duration"`
}

type bookingResponse struct {
	Err     error   `json:"err,omitempty"`
	Booking Booking `json:"booking"`
}

func (r bookingResponse) error() error { return r.Err }

type deleteRequest struct {
	BookingId string `json:"id"`
}

type deleteResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteResponse) error() error { return r.Err }

type getAllResponse struct {
	Err      error     `json:"err,omitempty"`
	Bookings []Booking `json:"bookings"`
}

func (r getAllResponse) error() error { return r.Err }
