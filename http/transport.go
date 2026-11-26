package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"example.com/car-service/car"
	"example.com/car-service/transport"
)

var errInvalidJSON = errors.New("invalid JSON")

// NewService связывает эндпоинты Go kit с HTTP-транспортом
func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/cars").Handler(kithttp.NewServer(
		endpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/cars/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PATCH").Path("/cars/{id}/sold").Handler(kithttp.NewServer(
		endpoints.MarkSold,
		decodeMarkSoldRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Car); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidJSON, err)
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.GetByIDRequest{ID: vars["id"]}, nil
}

func decodeMarkSoldRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.MarkSoldRequest{ID: vars["id"]}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, errInvalidJSON):
		status = http.StatusBadRequest
	case errors.Is(err, car.ErrCarNotFound):
		status = http.StatusNotFound
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
