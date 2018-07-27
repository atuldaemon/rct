package booking

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"bytes"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/booking/v1/").Handler(httptransport.NewServer(
		e.GetAllEndpoint,
		decodeGetAllRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/booking/v1/").Handler(httptransport.NewServer(
		e.BookingEndpoint,
		decodeBookingRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/booking/v1/{id}").Handler(httptransport.NewServer(
		e.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeGetAllRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getAllRequest
	return req, nil
}

func decodeBookingRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req bookingRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeBookingResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response bookingResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteRequest{BookingId: id}, nil
}

func decodeDeleteResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

func encodeBookingRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/booking/"
	return encodeRequest(ctx, req, request)
}

func encodeDeleteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/booking/"
	return encodeRequest(ctx, req, request)
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
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
