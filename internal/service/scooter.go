package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/integer00/e-scooter/internal/entity"
	"github.com/integer00/e-scooter/internal/registry"
	"github.com/integer00/e-scooter/internal/repository"
)

type scooterService struct {
	scooterRegistry   *registry.ScooterRegistry
	scooterRepository *repository.ScooterRepository
	paymentGate       entity.PaymentGateway
	postgresRepo      *repository.PostgresRepo
}

// add interfaces
func NewService(sr *registry.ScooterRegistry,
	screpo *repository.ScooterRepository, pg entity.PaymentGateway,
	pgr *repository.PostgresRepo) entity.ScooterService {
	return &scooterService{
		scooterRegistry:   sr,
		scooterRepository: screpo,
		paymentGate:       pg,
		postgresRepo:      pgr,
	}
}

func (suc scooterService) UserLogin(name string) (string, error) {
	log.Trace("usecase for userlogin")
	// ctx := context.Background() //timeout

	// log.Info(suc.postgresRepo.GetRides(context.Background()))

	if user := suc.scooterRegistry.GetUserById(name); user != nil {
		return name, nil
	} else { //creating user (from db)
		//check db
		log.Info("going for user to db")
		// query := "select name from users where name = '" + name + "'"
		// if dbUser, err := suc.postgresRepo.GetUserById(ctx, query); err == nil {
		// log.Info(dbUser)
		// log.Info("found user in db, adding to local registry")

		// err := suc.userRegistry.AddUser(dbUser)
		// if err != nil {
		// 	log.Error("failed to create user!")
		// 	return "", nil

		// } else {
		// log.Info("adding to db")
		// query := "insert into users (name) values ('" + name + "')"
		// err := suc.postgresRepo.AddUser(ctx, query)
		// if err != nil {
		// log.Error(err)
		// }

		er := suc.scooterRegistry.AddUser(name)
		if er != nil {
			log.Error("failed to create user!")
			return "", nil
		}
		// }
	}

	return name, nil
}

func (suc scooterService) GetScooter(s string) string {
	log.Trace("usecase for getting scooter")

	// suc.scooterRegistry.GetScooter()
	return "usecase returns"
}

func (suc scooterService) RegisterScooter(scooter *entity.Scooter) error {
	log.Trace("usecase for scooterRegistry")

	err := suc.scooterRegistry.RegisterScooter(*scooter)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (suc scooterService) GetEndpoints() []byte {
	log.Trace("usecase for getendpoints")
	// rides, _ := suc.postgresRepo.GetRides(context.Background())

	// log.Printf("%+v", rides)

	return suc.scooterRegistry.GetScooters()
}

func (suc scooterService) BookScooter(scooterId string, userId string) error {
	log.Info("booking scooter...")
	ctx := context.Background()

	sco, err := suc.scooterRegistry.GetScooterById(scooterId)
	if err != nil {
		return errors.New("failed to get scooter by that id")
	}

	if !sco.Available {
		return errors.New("requested scooter is not available")
	}

	user := suc.scooterRegistry.GetUserById(userId)
	if user == nil {
		log.Error("user not found")
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

	// suc.scooterRegistry.SetStatus(*sco, false)

	sco.Available = false

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

func (suc scooterService) StartScooter(scooterId string, userId string) error {
	log.Trace("usecase for starting scooter")
	//handle all related things , charge user, start scooter
	// validate and start has different context
	// ctx := context.Background()

	scooter, err := suc.scooterRegistry.GetScooterById(scooterId)
	if err != nil {
		return err
	}

	// rides, err := suc.postgresRepo.GetActiveRide(ctx)
	// if err != nil {
	// 	log.Error(err)
	// 	return errors.New("no rides been found")
	// }

	// if rides.UserId != userId {
	// 	return errors.New("asked from" + userId + " but booked for " + rides.UserId)
	// }

	// log.Info("got ride ", rides.RideId)

	// log.Printf("%+v", rides)

	suc.paymentGate.ChargeDeposit() //need action type(firstStart\finishRide)

	suc.scooterRepository.StartScooter(*scooter)

	// timeStart := time.Now().Unix()

	// rides.Status = "RIDE_IN_PROGRESS"
	// rides.StartTime = &timeStart

	// sql := fmt.Sprintf("update rides set status = '%s', start_time = '%d' where ride_id = '%s'",
	// 	rides.Status, *rides.StartTime, rides.RideId)
	// suc.postgresRepo.UpdateRide(ctx, sql)

	// log.Infof("%+v\n", rides)

	return nil
}

func (suc scooterService) StopScooter(scooterId string, userId string) error {
	log.Trace("usecase for stoping scooter")

	// ctx := context.Background()

	scooter, err := suc.scooterRegistry.GetScooterById(scooterId)
	if err != nil {
		return err
	}
	// rides, err := suc.postgresRepo.GetActiveRide(ctx)
	// if err != nil {
	// 	log.Error(err)
	// 	return errors.New("no rides been found")
	// }

	// if rides.UserId != userId {
	// 	return errors.New("asked from: " + userId + " but booked for: " + rides.UserId)
	// }

	// log.Printf("%+v", rides)

	suc.paymentGate.ChargeFair()
	// fare := 5 //get real fare

	suc.scooterRepository.StopScooter(*scooter)

	//releasing
	scooter.Available = true

	// timeStop := time.Now().Unix()

	// // rides.StopTime = stopTime
	// rides.Status = "DONE"
	// rides.StopTime = &timeStop
	// rides.FareCharged = &fare

	// sql := fmt.Sprintf("update rides set status = '%s', stop_time = '%d', fare_charged = '%d' where ride_id = '%s'",
	// 	rides.Status, *rides.StopTime, *rides.FareCharged, rides.RideId)

	// suc.postgresRepo.UpdateRide(ctx, sql)
	return nil
}

func (sco scooterService) RideHistory(userId string) {
	log.Trace("usecase for ride history")
	ctx := context.Background()

	rides, err := sco.postgresRepo.GetRidesById(ctx, userId)
	if err != nil {
	}

	log.Info(rides)

}

func (ss scooterService) GetUsers() {
	ss.scooterRegistry.GetUsers()
}
