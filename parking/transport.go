package parking

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting   = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrInvalidParam = errors.New("invalid param")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/parking/v1/getAll/").Handler(httptransport.NewServer(
		e.GetAllParkingEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/parking/v1/getFree/").Handler(httptransport.NewServer(
		e.GetFreeParkingEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/parking/v1/getReserved/").Handler(httptransport.NewServer(
		e.GetReservedParkingEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/parking/v1/search/").Handler(httptransport.NewServer(
		e.SearchParkingEndpoint,
		decodeSearchRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/parking/v1/find/{id}").Handler(httptransport.NewServer(
		e.FindByIdParkingEndpoint,
		decodeFindRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/parking/v1/").Handler(httptransport.NewServer(
		e.UpdateParkingEndpoint,
		decodeUpdateRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req updateParkingRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeSearchRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req searchParkingRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	switch req.Metric {
	case COST:
	case DIST:
		return req, nil
	default:
		return req, ErrInvalidParam
	}
	return req, nil
}

func decodeFindRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return findByIdParkingRequest{ID: id}, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getAllParkingRequest
	return req, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
