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
	GetEndpoints() []byte
	RegisterScooter(s Scooter) error
	UserLogin(s string) (string, error)
	ValidateJWT(s string) bool //remove
}

type Scooter struct {
	Id      string `json:"id" validate:"required"`
	Address string `json:"address"`
}
