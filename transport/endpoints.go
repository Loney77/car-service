package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/car-service/car"
)

// DTO для запросов и ответов
type CreateRequest struct {
	Car car.Car
}

type CreateResponse struct {
	ID string `json:"id"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Car car.Car `json:"car"`
}

type MarkSoldRequest struct {
	ID string `json:"id"`
}

type MarkSoldResponse struct{}

// Endpoints - набор эндпоинтов сервиса автомобилей
type Endpoints struct {
	Create   endpoint.Endpoint
	GetByID  endpoint.Endpoint
	MarkSold endpoint.Endpoint
}

// MakeEndpoints создаёт эндпоинты для сервиса автомобилей
func MakeEndpoints(s car.Service) Endpoints {
	return Endpoints{
		Create:   makeCreateEndpoint(s),
		GetByID:  makeGetByIDEndpoint(s),
		MarkSold: makeMarkSoldEndpoint(s),
	}
}

func makeCreateEndpoint(s car.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Car)
		return CreateResponse{ID: id}, err
	}
}

func makeGetByIDEndpoint(s car.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		c, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Car: c}, err
	}
}

func makeMarkSoldEndpoint(s car.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MarkSoldRequest)
		err := s.MarkSold(ctx, req.ID)
		return MarkSoldResponse{}, err
	}
}
