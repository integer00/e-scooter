package repository

import (
	log "github.com/sirupsen/logrus"
)

type User struct {
	Id   string
	Name string
}

type UserRegistry struct {
	userRegistry []User
}

func NewUserRegistry() *UserRegistry {
	return &UserRegistry{
		userRegistry: []User{},
	}
}

func (ur *UserRegistry) GetUser(s string) bool {

	for i := range ur.userRegistry {
		if s == ur.userRegistry[i].Name {
			log.Info("found user in registry")
			return true
		}
	}

	return false

}

func (ur *UserRegistry) AddUser(s string) error {
	ur.userRegistry = append(ur.userRegistry, User{
		Id:   "some",
		Name: s,
	})

	return nil
}
