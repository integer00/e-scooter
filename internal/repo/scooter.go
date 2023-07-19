package repo

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/integer00/e-scooter/internal/models"
	log "github.com/sirupsen/logrus"
)

type Scooter struct {
	ID      string `json:"id" validate:"required"`
	Address string `json:"address" validate:"required"`
}

var scooterRegistry = make(map[string]string)

func (u Scooter) Start() error {
	println("starting with" + u.Address)

	DoHTTPRequest("POST", []byte(u.ID), "http://"+u.Address+"/start")

	return nil
}

func (u Scooter) Stop() error {
	println("stopping with" + u.Address)

	DoHTTPRequest("POST", []byte(u.ID), "http://"+u.Address+"/stop")

	return nil
}

func NewScooter(id string, address string) models.Scooter {
	return &Scooter{
		ID:      id,
		Address: address,
	}
}

func startScooter(s Scooter) error {
	contactScooter(s, "start")
	println("scooter is started!")
	return nil
}

func stopScooter(s Scooter) error {
	contactScooter(s, "stop")
	println("scooter is stopped!")
	return nil
}

func contactScooter(s Scooter, action string) {
	var scooter = s

	switch action {
	case "start":
		scooter.Start()
	case "stop":
		scooter.Stop()
	default:
		log.Println("bad request")
	}
}

func CheckScooters() {
	println("checking scooters")

}

func scooterLookup(id string) (Scooter, error) {
	log.Println("asking for lookup:", id)

	for k, v := range scooterRegistry {
		if k == id {
			log.Println("found id in registry, " + id)

			return Scooter{ID: k, Address: v}, nil
		}
	}
	return Scooter{ID: "kappa", Address: "kappa"}, nil //?
}

func ParseScooter(req http.Request) Scooter {
	var sco Scooter
	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&sco)
	if err != nil {
		panic(err)
	}
	if err := validate.Struct(sco); err != nil {
		panic(err)
	}
	return sco
}

func RegisterScooter(s Scooter) {

	addToRegistry(s)
	log.Println("registred")
}

// func validateScooter(s string) Scooter {

// }

func addToRegistry(s Scooter) {
	//check if id is unique

	log.Println("adding to registry...")
	scooterRegistry[s.ID] = s.Address
	log.Println(scooterRegistry)
}

func ScooterHandler(s models.Message, action string) error {

	scooter, _ := scooterLookup(s.ID)

	// call that scooter and activate it
	switch action {
	case "startAction":
		startScooter(scooter)
		return nil
	default:
		stopScooter(scooter)
	}

	return nil
}
