package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"github.com/integer00/e-scooter/internal/entity"
	"github.com/integer00/e-scooter/internal/repository"
)

type scooterUseCase struct {
	scooterRegistry *repository.ScooterRegistry
	scooterApp      entity.ScooterService
	paymentGate     entity.PaymentGateway
	userRegistry    *repository.UserRegistry
}

func NewUseCase(sr *repository.ScooterRegistry,
	sapp entity.ScooterService, pg entity.PaymentGateway,
	ur *repository.UserRegistry) *scooterUseCase {
	return &scooterUseCase{
		scooterRegistry: sr,
		scooterApp:      sapp,
		paymentGate:     pg,
		userRegistry:    ur,
	}
}

var jwtkey = []byte("somesecretkey")

func (suc scooterUseCase) UserLogin(name string) (string, error) {
	log.Trace("usecase for userlogin")

	//somesort of userregistry is needed
	//check if user is there
	ok := suc.userRegistry.GetUser(name)
	if ok {
		s := suc.generateJWT(name)
		log.Info(s)
		//token handling (check if exist, if valid)
		//exit
		return s, nil
	}
	err := suc.userRegistry.AddUser(name)
	if err != nil {
		log.Error("failed to create user!")
		return "", nil
	}
	s := suc.generateJWT(name)
	log.Info(s)
	//token handling (check if exist, if valid)
	//exit
	return s, nil
}

func (suc scooterUseCase) generateJWT(name string) string {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["foo"] = "bar"
	claims["user"] = name

	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Error("failed to sign key")
		log.Error(err)
	}
	log.Info("signing key...")
	log.Info(tokenString)

	return tokenString
}
func (suc scooterUseCase) ValidateJWT(s string) bool {

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		log.Error("parse error")
		log.Error(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Info(claims)
		log.Info(claims["user"], claims["exp"], claims["bar"])
	} else {
		log.Error("error with claims")
		fmt.Println(err)
	}

	return true

}

func (suc scooterUseCase) GetScooter(s string) string {
	log.Trace("usecase for getting scooter")

	// suc.scooterRegistry.GetScooter()
	return "usecase returns"
}

func (suc scooterUseCase) RegisterScooter(scooter entity.Scooter) error {
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

	sc, err := suc.scooterRegistry.GetScooterById(ctx.Value("scooter").(string))
	if err != nil {
		return err
	}

	suc.paymentGate.ChargeDeposit() //need action type(firstStart\finishRide)

	suc.scooterApp.StartScooter(ctx, sc)
	return nil
}

func (suc scooterUseCase) StopScooter(ctx context.Context) error {
	log.Trace("usecase for stoping scooter")

	if ctx.Value("scooter") == nil {
		panic("no scooter id")
	}

	sc, err := suc.scooterRegistry.GetScooterById(ctx.Value("scooter").(string))
	if err != nil {
		return err
	}

	suc.paymentGate.ChargeFair()

	suc.scooterApp.StopScooter(ctx, sc)
	return nil
}

func (suc scooterUseCase) GetEndpoints() []byte {
	log.Trace("usecase for getendpoints")

	return suc.scooterRegistry.GetScooters()
}
