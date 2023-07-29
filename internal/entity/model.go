package entity

import (
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
	StartScooter(sc Scooter) error
	StopScooter(sc Scooter) error
}

type UseCase interface {
	GetScooter(s string) string
	StartScooter(s string) error
	StopScooter(s string) error
	GetEndpoints() []byte
	RegisterScooter(s *Scooter) error
	UserLogin(s string) (string, error)
	ValidateJWT(s string) bool //remove
	BookScooter(scooterID string, userID string) error
}

type User struct {
	Id   string
	Name string
}

type Scooter struct {
	Id      string `json:"id" validate:"required"`
	Address string `json:"address"`
}

type Message struct {
	Userid    string `json:"userid"`
	Scooterid string `json:"scooterid"`
}

type Ride struct {
	RideID  string  `db:"ride_id"`
	Scooter Scooter `db:"scooter_id"`
	User    User    `db:"user_id"`
	// Date        string
	// Time        string
	Status string `db:"status"`
	// FareCharged string
	// Distance    string
	// StartTime   string
	// StopTime    string
}
