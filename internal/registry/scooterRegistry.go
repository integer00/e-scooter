package registry

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

type ScooterRegistry struct {
	// registry map[string]string
	registry []entity.Scooter
	lock     *sync.Mutex
	// rideHistory []entity.Ride //implement cache in v2
}

// TODO need interface for ScooterRegistry
func NewRegistry() *ScooterRegistry {
	return &ScooterRegistry{
		registry: []entity.Scooter{},
		lock:     &sync.Mutex{},
		// rideHistory: []entity.Ride{},
	}
}

func (sr *ScooterRegistry) RegisterScooter(scooter entity.Scooter) error {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	log.Info("New registration!")
	log.Info(scooter)
	sr.registry = append(sr.registry, scooter)

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
		if sr.registry[i].Available {
			s = append(s, sr.registry[i].Id)
		}
	}

	a, _ := json.Marshal(&endpoints{Id: s})

	log.Info("json: ", string(a))

	return a
}
