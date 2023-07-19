package models

type Message struct {
	ID string `json:"id" validate:"required"`
}

type PaymentGateway interface {
	ChargeDeposit() error
	ChargeFair() error
	GetUserBalance() (int, error)
}

type Scooter interface {
	Start() error
	Stop() error
}
