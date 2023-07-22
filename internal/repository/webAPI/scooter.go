package webapi

import (
	"bytes"
	"context"
	"net/http"

	"github.com/integer00/e-scooter/internal/entity"
	log "github.com/sirupsen/logrus"
)

//put in another folder maybe

type ScooterService struct {
}

func NewScooterAPP() *ScooterService {
	return &ScooterService{}
}

// from client post /start&id=id
func (sapp ScooterService) StartScooter(ctx context.Context, sc entity.Scooter) error {
	log.Trace("starting scooter")
	log.Info("Starting scooter")
	sapp.start(sc)

	return nil
}
func (sapp ScooterService) StopScooter(ctx context.Context, sc entity.Scooter) error {
	log.Trace("stopping scooter")
	log.Info("Stopping scooter")
	sapp.stop(sc)

	return nil
}

// func (sapp ScooterService) contactScooter(action string) error {
// 	return nil
// }

func (sapp ScooterService) start(sc entity.Scooter) error {
	println("starting with" + sc.Address)

	DoHTTPRequest("POST", []byte(sc.ID), "http://"+sc.Address+"/start")

	return nil
}

func (sapp ScooterService) stop(sc entity.Scooter) error {
	println("stopping with" + sc.Address)

	DoHTTPRequest("POST", []byte(sc.ID), "http://"+sc.Address+"/stop")

	return nil
}

func DoHTTPRequest(method string, payload []byte, url string) http.Response {

	bodyReader := bytes.NewReader(payload)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		println("request failed")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return *res
}
