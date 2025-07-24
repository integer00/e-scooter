package repository

import (
	"bytes"
	"net/http"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

//http interface that interacts with scooters

type ScooterRepository struct {
}

func NewScooterRepository() *ScooterRepository {
	return &ScooterRepository{}
}

// from client post /start&id=id
func (sapp ScooterRepository) StartScooter(sc entity.Scooter) error {
	log.Trace("starting scooter")
	log.Info("Starting scooter")
	sapp.start(sc)

	return nil
}
func (sapp ScooterRepository) StopScooter(sc entity.Scooter) error {
	log.Trace("stopping scooter")
	log.Info("Stopping scooter")
	sapp.stop(sc)

	return nil
}

// func (sapp ScooterService) contactScooter(action string) error {
// 	return nil
// }

func (sapp ScooterRepository) start(sc entity.Scooter) error {
	println("starting with" + sc.Address)

	res := DoHTTPRequest("POST", []byte(sc.Id), "http://"+sc.Address+"/start")
	defer res.Body.Close()

	return nil
}

func (sapp ScooterRepository) stop(sc entity.Scooter) error {
	println("stopping with" + sc.Address)

	res := DoHTTPRequest("POST", []byte(sc.Id), "http://"+sc.Address+"/stop")
	defer res.Body.Close()

	return nil
}

func DoHTTPRequest(method string, payload []byte, url string) *http.Response {

	bodyReader := bytes.NewReader(payload)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		println("request failed")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return res
}
