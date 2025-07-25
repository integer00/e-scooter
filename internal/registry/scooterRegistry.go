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
	userRegistry    []entity.User
	scooterRegistry []entity.Scooter
	rideRegistry    []entity.Ride
	lock            *sync.Mutex
	// rideHistory []entity.Ride //implement cache in v2
}

// TODO need interface for ScooterRegistry
func NewRegistry() *ScooterRegistry {
	return &ScooterRegistry{
		userRegistry:    []entity.User{},
		scooterRegistry: []entity.Scooter{},
		rideRegistry:    []entity.Ride{},
		lock:            &sync.Mutex{},
	}
}

func (sr *ScooterRegistry) RegisterScooter(scooter entity.Scooter) error {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	log.Info("New registration!")
	log.Info(scooter)
	sr.scooterRegistry = append(sr.scooterRegistry, scooter)

	return nil
}

func (sr ScooterRegistry) GetScooterById(s string) (*entity.Scooter, error) {
	log.Println("asking for lookup:", s)

	for i := range sr.scooterRegistry {
		if s == sr.scooterRegistry[i].Id {
			log.Info("found match for id")
			return &sr.scooterRegistry[i], nil

		}
	}
	return nil, errors.New("error finding scooter")
}

type endpoints struct {
	Id []string
}

func (sr *ScooterRegistry) GetScooters() []byte {

	// {"id": ["a","b","c"]}

	log.Info("registry: ", sr.scooterRegistry)

	s := []string{}

	for i := range sr.scooterRegistry {
		if sr.scooterRegistry[i].Available {
			s = append(s, sr.scooterRegistry[i].Id)
		}
	}

	a, _ := json.Marshal(&endpoints{Id: s})

	log.Info("json: ", string(a))

	return a
}

func (sr *ScooterRegistry) GetUsers() {
	log.Info("user registry: ", sr.userRegistry)
}

func (sr *ScooterRegistry) GetUserById(s string) *entity.User {

	for i := range sr.userRegistry {
		if s == sr.userRegistry[i].Name {
			log.Info("found user in registry")
			return &sr.userRegistry[i]
		}
	}

	return nil
}

func (sr *ScooterRegistry) AddUser(s string) error {
	log.Info("adding user to registry")
	sr.lock.Lock()
	defer sr.lock.Unlock()
	log.Info(s)

	sr.userRegistry = append(sr.userRegistry, entity.User{
		// Id:   "some",
		Name: s,
	})

	return nil
}
