package entity

import (
	"context"
	"net/http"
)

type Controller interface {
	NewMux() *http.ServeMux
}

// type Registry interface {
// 	// GetScooter(sc interface{})
// 	RegisterScooter(r Registry, sc *Scooter) error
// }

type PaymentGateway interface {
	ChargeDeposit() error
	ChargeFair() error
	GetUserBalance() (int, error)
}

type ScooterService interface {
	StartScooter(ctx context.Context, sc Scooter) error
	StopScooter(ctx context.Context, sc Scooter) error
}

type UseCase interface {
	GetScooter(s string) string
	StartScooter(ctx context.Context) error
	StopScooter(ctx context.Context) error
	GetEndpoints() string
	RegisterScooter(s *Scooter) error
}

type Scooter struct {
	ID      string `json:"id" validate:"required`
	Address string `json:"address"`
}
