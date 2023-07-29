package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/integer00/e-scooter/internal/entity"
	"github.com/integer00/e-scooter/internal/repository"
)

type scooterUseCase struct {
	scooterRegistry *repository.ScooterRegistry
	scooterApp      entity.ScooterService
	paymentGate     entity.PaymentGateway
	userRegistry    *repository.UserRegistry
	postgresRepo    *repository.PostgresRepo
}

func NewUseCase(sr *repository.ScooterRegistry,
	sapp entity.ScooterService, pg entity.PaymentGateway,
	ur *repository.UserRegistry, pgr *repository.PostgresRepo) *scooterUseCase {
	return &scooterUseCase{
		scooterRegistry: sr,
		scooterApp:      sapp,
		paymentGate:     pg,
		userRegistry:    ur,
		postgresRepo:    pgr,
	}
}

var jwtkey = []byte("somesecretkey")

func (suc scooterUseCase) UserLogin(name string) (string, error) {
	log.Trace("usecase for userlogin")
	var userToSign = name
	ctx := context.Background() //timeout

	log.Info(suc.postgresRepo.GetUsers(ctx))
	log.Info(suc.postgresRepo.GetRides(ctx))

	if user := suc.userRegistry.GetUserById(userToSign); user != nil {
		s := suc.generateJWT(userToSign)

		//token handling (check if exist, if valid)
		//exit
		return s, nil
	} else { //creating user (from db)
		//check db
		log.Info("going for user to db")
		query := "select name from users where name = '" + name + "'"
		if dbUser, err := suc.postgresRepo.FindUserById(ctx, query); err == nil {
			log.Info(dbUser)
			log.Info("found user in db, adding to local registry")

			err := suc.userRegistry.AddUser(dbUser)
			if err != nil {
				log.Error("failed to create user!")
				return "", nil
			}

		} else {
			log.Info("adding to db")
			query := "insert into users (name) values ('" + name + "')"
			err := suc.postgresRepo.AddUserById(ctx, query)
			if err != nil {
				log.Error(err)
			}

			er := suc.userRegistry.AddUser(name)
			if er != nil {
				log.Error("failed to create user!")
				return "", nil
			}
		}
	}

	s := suc.generateJWT(userToSign)
	// log.Info(s)
	return s, nil
	//token handling (check if exist, if valid)
	//exit
}

func (suc scooterUseCase) generateJWT(name string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":  time.Now().Add(60 * time.Minute).Unix(),
		"user": name,
	})
	//todo proper claims map

	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Error("failed to sign key")
		log.Error(err)
	}
	log.Info("signing key...")

	return tokenString
}
func (suc scooterUseCase) ValidateJWT(s string) bool {
	//validate token and claims
	//check expiration

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		log.Error("parse error")
		log.Error(err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Info(claims)
		log.Info(claims["user"], claims["exp"])
	} else {
		log.Error("error with claims")
		fmt.Println(err)
		return false
	}

	return true

}

func (suc scooterUseCase) GetScooter(s string) string {
	log.Trace("usecase for getting scooter")

	// suc.scooterRegistry.GetScooter()
	return "usecase returns"
}

func (suc scooterUseCase) RegisterScooter(scooter *entity.Scooter) error {
	log.Trace("usecase for scooterRegistry")

	suc.scooterRegistry.RegisterScooter(*scooter)

	return nil
}

func (suc scooterUseCase) StartScooter(s string) error {
	log.Trace("usecase for starting scooter")
	//handle all related things , charge user, start scooter
	// validate and start has different context

	sc, err := suc.scooterRegistry.GetScooterById(s)
	if err != nil {
		return err
	}

	suc.paymentGate.ChargeDeposit() //need action type(firstStart\finishRide)

	suc.scooterApp.StartScooter(*sc)
	return nil
}

func (suc scooterUseCase) StopScooter(s string) error {
	log.Trace("usecase for stoping scooter")

	sc, err := suc.scooterRegistry.GetScooterById(s)
	if err != nil {
		return err
	}

	suc.paymentGate.ChargeFair()

	suc.scooterApp.StopScooter(*sc)
	return nil
}

func (suc scooterUseCase) GetEndpoints() []byte {
	log.Trace("usecase for getendpoints")

	return suc.scooterRegistry.GetScooters()
}

func (suc scooterUseCase) BookScooter(scooterID string, userID string) error {
	log.Info("booking scooter...")
	ctx := context.Background()

	sco, err := suc.scooterRegistry.GetScooterById(scooterID)
	if err != nil {
		return errors.New("failed to get scooter by that id")
	}
	user := suc.userRegistry.GetUserById(userID)
	if user == nil {
		return errors.New("failed to get user by that id")
	}

	uuid := uuid.New()

	//booking scooter, making it belong to userid
	//doing query to db
	//adding record to rideHistory

	ride := entity.Ride{
		RideID:  uuid.String(),
		Scooter: *sco,
		User:    *user,
		Status:  "booking",
	}

	//add ride to ride history
	if err := suc.scooterRegistry.AddRide(ride); err != nil {
		log.Error("failed to add ride to histry")
		log.Error(err)
	}

	//add to db
	if err := suc.postgresRepo.AddRide(ctx, ride); err != nil {
		log.Error("failed to add ride to db")
		log.Error(err)
	}

	return nil
}
