package car

import "context"

type Car struct {
	ID    string  `json:"id,omitempty"`
	Brand string  `json:"brand"`
	Model string  `json:"model"`
	Year  int     `json:"year"`
	Price float64 `json:"price"`
	Sold  bool    `json:"sold"`
}

// Service описывает бизнес-логику сервиса автомобилей
type Service interface {
	Create(ctx context.Context, car Car) (string, error)
	GetByID(ctx context.Context, id string) (Car, error)
	MarkSold(ctx context.Context, id string) error
}

// Repository описывает хранилище автомобилей
type Repository interface {
	CreateCar(ctx context.Context, car Car) error
	GetCarByID(ctx context.Context, id string) (Car, error)
	MarkCarSold(ctx context.Context, id string) error
}
