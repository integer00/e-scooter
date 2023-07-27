package repository

import (
	"errors"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

type UserRegistry struct {
	userRegistry []entity.User
}

func NewUserRegistry() *UserRegistry {
	return &UserRegistry{
		userRegistry: []entity.User{},
	}
}

func (ur *UserRegistry) GetUserById(s string) (*entity.User, error) {

	for i := range ur.userRegistry {
		if s == ur.userRegistry[i].Name {
			log.Info("found user in registry")
			return &ur.userRegistry[i], nil
		}
	}

	return nil, errors.New("user is not found")

}

func (ur *UserRegistry) AddUser(s string) error {
	ur.userRegistry = append(ur.userRegistry, entity.User{
		// Id:   "some",
		Name: s,
	})

	return nil
}
