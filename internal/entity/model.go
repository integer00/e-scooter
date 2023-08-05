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
	StartScooter(sc Scooter) error
	StopScooter(sc Scooter) error
}

type User struct {
	Id   string
	Name string `json:"userid"`
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
