package parking

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetAllParkingEndpoint      endpoint.Endpoint
	GetFreeParkingEndpoint     endpoint.Endpoint
	GetReservedParkingEndpoint endpoint.Endpoint
	SearchParkingEndpoint      endpoint.Endpoint
	FindByIdParkingEndpoint    endpoint.Endpoint
	UpdateParkingEndpoint      endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetAllParkingEndpoint:      MakeGetAllEndpoint(s),
		GetFreeParkingEndpoint:     MakeGetFreeEndpoint(s),
		GetReservedParkingEndpoint: MakeGetReservedEndpoint(s),
		SearchParkingEndpoint:      MakeSearchEndpoint(s),
		FindByIdParkingEndpoint:    MakeFindByIdEndpoint(s),
		UpdateParkingEndpoint:      MakeUpdateEndpoint(s),
	}
}

func (e Endpoints) GetAllParking(ctx context.Context) ([]Spot, error) {
	request := getAllParkingRequest{}
	response, err := e.GetAllParkingEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(getAllParkingResponse)
	return resp.Spots, resp.Err
}

func (e Endpoints) GetFreeParking(ctx context.Context) ([]Spot, error) {
	request := getFreeParkingRequest{}
	response, err := e.GetFreeParkingEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(getFreeParkingResponse)
	return resp.Spots, resp.Err
}

func (e Endpoints) GetReservedParking(ctx context.Context) ([]Spot, error) {
	request := getReservedParkingRequest{}
	response, err := e.GetReservedParkingEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(getReservedParkingResponse)
	return resp.Spots, resp.Err
}

func MakeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(getAllParkingRequest)
		ss, e := s.GetAll(ctx)
		return getAllParkingResponse{Spots: ss, Err: e}, e
	}
}

func MakeGetFreeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(getAllParkingRequest)
		ss, e := s.GetFree(ctx)
		return getFreeParkingResponse{Spots: ss, Err: e}, e
	}
}

func MakeGetReservedEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(getAllParkingRequest)
		ss, e := s.GetReserved(ctx)
		return getReservedParkingResponse{Spots: ss, Err: e}, e
	}
}

func MakeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateParkingRequest)
		s, err := s.Update(ctx, req.Spot)
		return updateParkingResponse{Spot: s, Err: err}, err
	}
}

func MakeSearchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(searchParkingRequest)
		ss, e := s.Search(ctx, req.Lat, req.Lon, req.Rad, req.Metric)
		return getSearchParkingResponse{Spots: ss, Err: e}, e
	}
}

func MakeFindByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findByIdParkingRequest)
		s, e := s.FindById(ctx, req.ID)
		return getAllParkingResponse{Spots: []Spot{s}, Err: e}, e
	}
}

//

type updateParkingResponse struct {
	Err  error `json:"err,omitempty"`
	Spot Spot  `json:"spots"`
}

func (r updateParkingResponse) error() error { return r.Err }

type findByIdParkingRequest struct {
	ID string `json:"id"`
}
type updateParkingRequest struct {
	Spot Spot `json:"spot"`
}

type SearchMetric string

const (
	COST  SearchMetric = "cost"
	DIST  SearchMetric = "dist"
)

type searchParkingRequest struct {
	Lat    string       `json:"lat"`
	Lon    string       `json:"lon"`
	Rad    string       `json:"rad"`
	Metric SearchMetric `json:"metric"`
}

type getAllParkingRequest struct {
}

type getParkingResponse struct {
	Err   error  `json:"err,omitempty"`
	Spots []Spot `json:"spots"`
}

func (r getParkingResponse) error() error { return r.Err }

type getAllParkingResponse struct {
	Err   error  `json:"err,omitempty"`
	Spots []Spot `json:"spots"`
}

func (r getAllParkingResponse) error() error { return r.Err }

type getSearchParkingResponse struct {
	Err   error          `json:"err,omitempty"`
	Spots []ExtendedSpot `json:"spots"`
}

func (r getSearchParkingResponse) error() error { return r.Err }

type getFreeParkingRequest struct {
}

type getFreeParkingResponse struct {
	Err   error  `json:"err,omitempty"`
	Spots []Spot `json:"spots"`
}

func (r getFreeParkingResponse) error() error { return r.Err }

type getReservedParkingRequest struct {
}

type getReservedParkingResponse struct {
	Err   error  `json:"err,omitempty"`
	Spots []Spot `json:"spots"`
}

func (r getReservedParkingResponse) error() error { return r.Err }
