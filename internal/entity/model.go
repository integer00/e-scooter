package entity

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
	BookScooter(scooterId string, userId string) error
	StartScooter(scooterId string, userId string) error
	StopScooter(scooterId string, userId string) error
	RideHistory(userId string)
	GetScooter(s string) string
	GetEndpoints() []byte
	RegisterScooter(s *Scooter) error
	UserLogin(s string) (string, error)
	GetUsers()
}

type User struct {
	Id   string
	Name string
	// Name string `json:"userid"`
}

type Scooter struct {
	Id        string `json:"id" validate:"required"`
	Address   string `json:"address"`
	Available bool   `json:"available"`
}

type Message struct {
	UserId    string `json:"userid" validate:"required"`
	ScooterId string `json:"scooterid" validate:"required"`
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
