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
}

// TODO need interface for ScooterRegistry
func NewRegistry() *ScooterRegistry {
	return &ScooterRegistry{
		registry: []entity.Scooter{},
	}
}

//	func GetScooter(i interface{}) *entity.Scooter {
//		log.Trace("getting scooter")
//		return &entity.Scooter{
//			ID:      i.id,
//			Address: i.address,
//		}
//	}
// func RegisterScooter(sc *entity.Scooter) error {
// 	return nil

// }

func (sr *ScooterRegistry) RegisterScooter(scooter entity.Scooter) error {
	log.Info("New registration!")
	sr.registerScooter(scooter)

	return nil
}

func (sr ScooterRegistry) GetScooterById(s string) (entity.Scooter, error) {

	log.Info(sr.registry)

	sc, err := sr.scooterLookupById(s)
	if err != nil {
		return entity.Scooter{}, errors.New("Scooter is not found!")
	}

	return sc, nil
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

// func (sr ScooterRegistry) newScooter(id string, address string) entity.Scooter {
// 	return entity.Scooter{
// 		Id:      id,
// 		Address: address,
// 	}
// }

// func (sr ScooterRegistry) getScooter(scooter *entity.Scooter) {
// 	// return sr.registry
// 	// return "here bro"
// }

func (sr *ScooterRegistry) registerScooter(scooter entity.Scooter) {
	log.Println("adding to registry...")

	sr.registry = append(sr.registry, scooter)

	log.Info(sr.registry)
}

func (sr ScooterRegistry) scooterLookupById(id string) (entity.Scooter, error) {
	log.Println("asking for lookup:", id)

	for i := range sr.registry {
		if id == sr.registry[i].Id {
			log.Info("found match for id")
			return sr.registry[i], nil

		}
	}
	return entity.Scooter{Id: "kappa", Address: "kappa"}, errors.New("error finding scooter") //?
}
