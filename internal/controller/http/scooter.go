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
	mux.HandleFunc("/endpoints", sc.getEndpoints)
	mux.HandleFunc("/start", sc.startScooterHandler)
	mux.HandleFunc("/stop", sc.stopScooterHandler)
	return mux
}

func (sc ScoController) registerEndpointHandler(w http.ResponseWriter, req *http.Request) {
	p := parseRequest(*req)

	sc.scooterUseCase.RegisterScooter(p)
}

func (sc ScoController) getEndpoints(w http.ResponseWriter, req *http.Request) {
	log.Println("endpoints at controller, sending to usecase")

	// s := sc.scooterUseCase.GetEndpoints()

	// io.WriteString(w, s)

	// var response = scooterRegistry.(string)

	// log.Fprintf(w, scooterRegistry)
}

// might be also /api/scooter/:id/(start\stop)
func (sc ScoController) startScooterHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for start")

	s := parseRequest(*req)

	ctx := context.WithValue(context.Background(), "scooter", s.ID)

	sc.scooterUseCase.StartScooter(ctx)
}
func (sc ScoController) stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for start")

	s := parseRequest(*req)

	ctx := context.WithValue(context.Background(), "scooter", s.ID)

	sc.scooterUseCase.StopScooter(ctx)
}

func parseRequest(req http.Request) *entity.Scooter {
	var s = &entity.Scooter{}

	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&s)
	if err != nil {
		log.Error("could not decode json")
		panic(err)

	}
	if err := validate.Struct(s); err != nil {
		log.Error("could not validate json")
		panic(err)
	}
	return s
}
