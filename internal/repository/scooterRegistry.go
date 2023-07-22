package repository

import (
	"errors"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

type ScooterRegistry struct {
	registry map[string]string
}

// TODO need interface for ScooterRegistry
func NewRegistry() *ScooterRegistry {
	return &ScooterRegistry{
		registry: make(map[string]string),
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

func (sr ScooterRegistry) RegisterScooter(scooter *entity.Scooter) error {
	log.Info("New registration!")
	sr.registerScooter(scooter)

	return nil
}

func (sr ScooterRegistry) GetScooterById(s string) entity.Scooter {

	sc, err := sr.scooterLookupById(s)
	if err != nil {
		panic(err)
	}

	return sc
}

func (sr ScooterRegistry) newScooter(id string, address string) entity.Scooter {
	return entity.Scooter{
		ID:      id,
		Address: address,
	}
}

func (sr ScooterRegistry) getScooter(scooter *entity.Scooter) {
	// return sr.registry
	// return "here bro"
}

func (sr ScooterRegistry) registerScooter(scooter *entity.Scooter) {

	log.Println("adding to registry...")
	sr.registry[scooter.ID] = scooter.Address
}

func (sr ScooterRegistry) scooterLookupById(id string) (entity.Scooter, error) {
	log.Println("asking for lookup:", id)

	for k, v := range sr.registry {
		if k == id {
			log.Println("found id in registry, " + id)

			return entity.Scooter{ID: k, Address: v}, nil
		}
	}
	return entity.Scooter{ID: "kappa", Address: "kappa"}, errors.New("error finding scooter") //?
}
