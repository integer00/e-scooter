package registry

import (
	"sync"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

type UserRegistry struct {
	userRegistry []entity.User
	lock         *sync.Mutex
}

func NewUserRegistry() *UserRegistry {
	return &UserRegistry{
		userRegistry: []entity.User{},
		lock:         &sync.Mutex{},
	}
}

func (ur *UserRegistry) GetUsers() {
	log.Info("user registry: ", ur.userRegistry)
}

func (ur *UserRegistry) GetUserById(s string) *entity.User {

	for i := range ur.userRegistry {
		if s == ur.userRegistry[i].Name {
			log.Info("found user in registry")
			return &ur.userRegistry[i]
		}
	}

	return nil

}

func (ur *UserRegistry) AddUser(s string) error {
	log.Info("adding user to registry")
	ur.lock.Lock()
	defer ur.lock.Unlock()
	log.Info(s)

	ur.userRegistry = append(ur.userRegistry, entity.User{
		// Id:   "some",
		Name: s,
	})

	return nil
}
