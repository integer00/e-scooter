package repository

import (
	"errors"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

type PaymentGate struct {
	//PG is an external service, mocking
}

func NewPG() entity.PaymentGateway {
	return &PaymentGate{}
}

func (pg PaymentGate) ChargeDeposit() error {
	log.Println("charging deposit")
	return nil
}

func (pg PaymentGate) ChargeFair() error {
	log.Println("charging fair")
	return nil
}
func (pg PaymentGate) GetUserBalance() (int, error) {
	return 1000, nil
}

func PgHandler(userid string, operation string) error {

	pg := NewPG()

	balance, _ := pg.GetUserBalance()

	switch operation {
	case "startRide":
		// check user balance
		if balance < 10 {
			log.Errorln("not enough balance!")
			return errors.New("error")
		}

		//charging is safe way
		pg.ChargeDeposit()
	case "endRide":
		var fair = 20

		if fair > balance {
			log.Println("returning deposit-fee")
			return nil
		}
		log.Println("returning deposit")
		pg.ChargeFair()
	default:
		log.Println("unknown operation")
	}

	return nil
}
