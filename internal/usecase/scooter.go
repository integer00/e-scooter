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

type UseCase interface {
	BookScooter(scooterId string, userId string) error
	StartScooter(scooterId string, userId string) error
	StopScooter(scooterId string, userId string) error
	RideHistory(userId string)
	GetScooter(s string) string
	GetEndpoints() []byte
	RegisterScooter(s *entity.Scooter) error
	UserLogin(s string) (string, error)
	ValidateJWT(s string) (jwt.MapClaims, bool) //remove
}

type scooterUseCase struct {
	scooterRegistry *repository.ScooterRegistry
	scooterApp      entity.ScooterService
	paymentGate     entity.PaymentGateway
	userRegistry    *repository.UserRegistry
	postgresRepo    *repository.PostgresRepo
}

func NewUseCase(sr *repository.ScooterRegistry,
	sapp entity.ScooterService, pg entity.PaymentGateway,
	ur *repository.UserRegistry, pgr *repository.PostgresRepo) UseCase {
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

	log.Info(suc.postgresRepo.GetRides(context.Background()))

	if user := suc.userRegistry.GetUserById(userToSign); user != nil {
		s := suc.generateJWT(userToSign)

		//token handling (check if exist, if valid)
		//exit
		return s, nil
	} else { //creating user (from db)
		//check db
		log.Info("going for user to db")
		query := "select name from users where name = '" + name + "'"
		if dbUser, err := suc.postgresRepo.GetUserById(ctx, query); err == nil {
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
			err := suc.postgresRepo.AddUser(ctx, query)
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
func (suc scooterUseCase) ValidateJWT(s string) (jwt.MapClaims, bool) {
	//validate token and claims
	//check expiration

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		log.Error("parse error")
		log.Error(err)
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		log.Info(claims)
		log.Info(claims["user"], claims["exp"])
	} else {
		log.Error("error with claims")
		fmt.Println(err)
		return nil, false
	}

	return claims, true

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

func (suc scooterUseCase) GetEndpoints() []byte {
	log.Trace("usecase for getendpoints")
	rides, _ := suc.postgresRepo.GetRides(context.Background())

	log.Printf("%+v", rides)

	return suc.scooterRegistry.GetScooters()
}

func (suc scooterUseCase) BookScooter(scooterId string, userId string) error {
	log.Info("booking scooter...")
	ctx := context.Background()

	sco, err := suc.scooterRegistry.GetScooterById(scooterId)
	if err != nil {
		return errors.New("failed to get scooter by that id")
	}
	user := suc.userRegistry.GetUserById(userId)
	if user == nil {
		return errors.New("failed to get user by that id")
	}

	//check if scooterId is already booked

	uuid := uuid.New()

	//booking scooter, making it belong to userid
	//doing query to db
	//adding record to rideHistory

	ride := entity.Ride{
		RideId:    uuid.String(),
		ScooterId: sco.Id,
		UserId:    user.Name,
		Status:    "booking",
	}
	log.Printf("%+v", ride)

	//add ride to ride history
	// if err := suc.scooterRegistry.AddRide(ride); err != nil {
	// 	log.Error("failed to add ride to histry")
	// 	log.Error(err)
	// }

	//add to db
	if err := suc.postgresRepo.AddRide(ctx, ride); err != nil {
		log.Error("failed to add ride to db")
		log.Error(err)
	}
	//scooter is gone from available pool

	return nil
}

func (suc scooterUseCase) StartScooter(scooterId string, userId string) error {
	log.Trace("usecase for starting scooter")
	//handle all related things , charge user, start scooter
	// validate and start has different context
	ctx := context.Background()

	scooter, err := suc.scooterRegistry.GetScooterById(scooterId)
	if err != nil {
		return err
	}

	rides, err := suc.postgresRepo.GetActiveRide(ctx)
	if err != nil {
		log.Error(err)
		return errors.New("no rides been founded")
	}
	log.Info("got ride ", rides.RideId)

	log.Printf("%+v", rides)

	suc.paymentGate.ChargeDeposit() //need action type(firstStart\finishRide)

	suc.scooterApp.StartScooter(*scooter)

	timeStart := time.Now().Unix()

	rides.Status = "RIDE_IN_PROGRESS"
	rides.StartTime = &timeStart

	sql := fmt.Sprintf("update rides set status = '%s', start_time = '%d' where ride_id = '%s'",
		rides.Status, *rides.StartTime, rides.RideId)
	suc.postgresRepo.UpdateRide(ctx, sql)

	log.Infof("%+v\n", rides)

	return nil
}

func (suc scooterUseCase) StopScooter(scooterId string, userId string) error {
	log.Trace("usecase for stoping scooter")

	ctx := context.Background()

	scooter, err := suc.scooterRegistry.GetScooterById(scooterId)
	if err != nil {
		return err
	}
	rides, err := suc.postgresRepo.GetActiveRide(ctx)
	if err != nil {
		log.Error(err)
		return errors.New("no rides been founded")
	}

	log.Printf("%+v", rides)

	suc.paymentGate.ChargeFair()
	fare := 5 //get real fare

	suc.scooterApp.StopScooter(*scooter)

	timeStop := time.Now().Unix()

	// rides.StopTime = stopTime
	rides.Status = "DONE"
	rides.StopTime = &timeStop
	rides.FareCharged = &fare

	sql := fmt.Sprintf("update rides set status = '%s', stop_time = '%d', fare_charged = '%d' where ride_id = '%s'",
		rides.Status, *rides.StopTime, *rides.FareCharged, rides.RideId)

	suc.postgresRepo.UpdateRide(ctx, sql)
	return nil
}

func (sco scooterUseCase) RideHistory(userId string) {
	log.Trace("usecase for ride history")
	ctx := context.Background()

	rides, err := sco.postgresRepo.GetRidesById(ctx, userId)
	if err != nil {
	}

	log.Info(rides)

}
