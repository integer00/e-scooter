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
	BookScooter(scooterID string, userID string) error
	StartScooter(scooterID string, userID string) error
	StopScooter(scooterID string, userID string) error
	GetScooter(s string) string
	GetEndpoints() []byte
	RegisterScooter(s *Scooter) error
	UserLogin(s string) (string, error)
	ValidateJWT(s string) bool //remove
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
	UserId    string `json:"userid"`
	ScooterId string `json:"scooterid"`
}

type Ride struct {
	RideId    string `db:"ride_id"`
	ScooterId string `db:"scooter_id"`
	UserId    string `db:"user_id"`
	Status    string `db:"status"`
	StartTime int64  `db:"start_time"`
	StopTime  int64  `db:"stop_time"`
	// Date        string
	// Time        string
	// FareCharged string
	// Distance    string
}
