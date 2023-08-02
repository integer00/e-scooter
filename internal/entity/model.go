package entity

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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
	BookScooter(scooterId string, userId string) error
	StartScooter(scooterId string, userId string) error
	StopScooter(scooterId string, userId string) error
	RideHistory(userId string)
	GetScooter(s string) string
	GetEndpoints() []byte
	RegisterScooter(s *Scooter) error
	UserLogin(s string) (string, error)
	ValidateJWT(s string) (jwt.MapClaims, bool) //remove
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
	RideId      string `db:"ride_id"`
	ScooterId   string `db:"scooter_id"`
	UserId      string `db:"user_id"`
	Status      string `db:"status"`
	StartTime   *int64 `db:"start_time"`
	StopTime    *int64 `db:"stop_time"`
	FareCharged *int   `db:"fare_charged"`
	// Distance    string
}
