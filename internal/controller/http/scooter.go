package http

import (
	"context"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator"
	"github.com/integer00/e-scooter/internal/entity"
)

type ScoController struct {
	scooterUseCase entity.UseCase
}

func NewScooterController(u entity.UseCase) entity.Controller {
	return &ScoController{
		scooterUseCase: u,
	}
}

func (sc ScoController) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", sc.registerEndpointHandler)
	mux.HandleFunc("/scooters", sc.getScootersHandler)
	mux.HandleFunc("/start", sc.startScooterHandler)
	mux.HandleFunc("/stop", sc.stopScooterHandler)
	return mux
}

func (sc ScoController) registerEndpointHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for registration")
	p := parseRequest(*req)

	sc.scooterUseCase.RegisterScooter(p)

	w.WriteHeader(http.StatusOK)

}

func (sc ScoController) getScootersHandler(w http.ResponseWriter, req *http.Request) {
	//this returns scooterID+geoCoordinates like {"id":"kappa_ride","location":"coordinates"}

	log.Info("asking for endpoints")

	s := sc.scooterUseCase.GetEndpoints()
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(s)

}

// might be also /api/scooter/:id/(start\stop)
func (sc ScoController) startScooterHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for start")

	s := parseRequest(*req)

	ctx := context.WithValue(context.Background(), "scooter", s.Id)

	err := sc.scooterUseCase.StartScooter(ctx)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

}
func (sc ScoController) stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for stop")

	s := parseRequest(*req)

	ctx := context.WithValue(context.Background(), "scooter", s.Id)

	err := sc.scooterUseCase.StopScooter(ctx)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func parseRequest(req http.Request) entity.Scooter {
	var s = entity.Scooter{}

	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&s)
	if err != nil {
		log.Warn(err)

	}
	if err := validate.Struct(s); err != nil {
		log.Warn(err)
	}
	return s
}
