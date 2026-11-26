package implementation

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"

	"example.com/car-service/car"
)

type service struct {
	repository car.Repository
	logger     log.Logger
}

// NewService создаёт новый экземпляр сервиса автомобилей
func NewService(repo car.Repository, logger log.Logger) car.Service {
	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, c car.Car) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	id := uuid.New().String()
	c.ID = id
	c.Sold = false

	if err := s.repository.CreateCar(ctx, c); err != nil {
		level.Error(logger).Log("err", err)
		return "", car.ErrCmdRepository
	}
	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (car.Car, error) {
	logger := log.With(s.logger, "method", "GetByID")

	c, err := s.repository.GetCarByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return c, car.ErrCarNotFound
	}
	return c, nil
}

func (s *service) MarkSold(ctx context.Context, id string) error {
	logger := log.With(s.logger, "method", "MarkSold")

	if err := s.repository.MarkCarSold(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}
	return nil
}
