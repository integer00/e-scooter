package repository

import (
	"encoding/json"
	"errors"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

type ScooterRegistry struct {
	// registry map[string]string
	registry []entity.Scooter
	// rideHistory []entity.Ride //implement cache in v2
}

// TODO need interface for ScooterRegistry
func NewRegistry() *ScooterRegistry {
	return &ScooterRegistry{
		registry: []entity.Scooter{},
		// rideHistory: []entity.Ride{},
	}
}

// func (sr *ScooterRegistry) AddRide(ride entity.Ride) error {
// 	log.Info("adding ride to rideHistory")
// 	sr.rideHistory = append(sr.rideHistory, ride)
// 	log.Info(sr.rideHistory)

// 	return nil
// }

// func (sr *ScooterRegistry) GetRideById(id string) error {
// 	log.Info("adding ride to rideHistory")
// 	sr.rideHistory = append(sr.rideHistory, ride)
// 	log.Info(sr.rideHistory)

// 	return nil
// }

func (sr *ScooterRegistry) RegisterScooter(scooter entity.Scooter) error {
	log.Info("New registration!")
	sr.registerScooter(scooter)

	return nil
}

func (sr ScooterRegistry) GetScooterById(s string) (*entity.Scooter, error) {
	log.Println("asking for lookup:", s)

	for i := range sr.registry {
		if s == sr.registry[i].Id {
			log.Info("found match for id")
			return &sr.registry[i], nil

		}
	}
	return nil, errors.New("error finding scooter")
}

type endpoints struct {
	Id []string
}

func (sr *ScooterRegistry) GetScooters() []byte {

	// {"id": ["a","b","c"]}

	log.Info("registry: ", sr.registry)

	s := []string{}

	for i := range sr.registry {
		s = append(s, sr.registry[i].Id)
	}

	a, _ := json.Marshal(&endpoints{Id: s})

	log.Info("json: ", string(a))

	return a
}

func (sr *ScooterRegistry) registerScooter(scooter entity.Scooter) {
	log.Println("adding to registry...")

	sr.registry = append(sr.registry, scooter)

	log.Info(sr.registry)
}
