package usecase

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/integer00/e-scooter/internal/entity"
	"github.com/integer00/e-scooter/internal/repository"
)

type scooterUseCase struct {
	scooterRegistry *repository.ScooterRegistry
	scooterApp      entity.ScooterService
	paymentGate     entity.PaymentGateway
}

func NewUseCase(sr *repository.ScooterRegistry, sapp entity.ScooterService, pg entity.PaymentGateway) *scooterUseCase {
	return &scooterUseCase{
		scooterRegistry: sr,
		scooterApp:      sapp,
		paymentGate:     pg,
	}
}

// func getScooter(scu scooterUseCase) string {
// 	return scu.scooterRegistry.GetScooter()
// }

func (suc scooterUseCase) GetScooter(s string) string {
	log.Trace("usecase for getting scooter")

	// suc.scooterRegistry.GetScooter()
	return "usecase returns"
}

func (suc scooterUseCase) RegisterScooter(scooter *entity.Scooter) error {
	log.Trace("usecase for scooterRegistry")

	suc.scooterRegistry.RegisterScooter(scooter)

	return nil
}

func (suc scooterUseCase) StartScooter(ctx context.Context) error {
	log.Trace("usecase for starting scooter")
	//handle all related things , charge user, start scooter
	// validate and start has different context

	if ctx.Value("scooter") == nil {
		panic("no scooter id")
	}

	sc := suc.scooterRegistry.GetScooterById(ctx.Value("scooter").(string))

	suc.paymentGate.ChargeDeposit() //need action type(firstStart\finishRide)

	suc.scooterApp.StartScooter(ctx, sc)
	return nil
}

func (suc scooterUseCase) StopScooter(ctx context.Context) error {
	log.Trace("usecase for stoping scooter")

	if ctx.Value("scooter") == nil {
		panic("no scooter id")
	}

	sc := suc.scooterRegistry.GetScooterById(ctx.Value("scooter").(string))

	suc.paymentGate.ChargeFair()

	suc.scooterApp.StopScooter(ctx, sc)
	return nil
}

func (suc scooterUseCase) GetEndpoints() string {
	log.Trace("at usecase getendpoints")
	//returning all
	return "here bro some endpoints from usecase"
}
