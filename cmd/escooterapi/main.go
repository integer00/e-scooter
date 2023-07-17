package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

type Scooter struct {
	ID      string `json:"id" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type Message struct {
	ID string `json:"id" validate:"required"`
}

func (u Scooter) Start() error {
	println("starting with" + u.Address)

	doHTTPRequest("POST", []byte(u.ID), "http://"+u.Address+"/start")

	return nil
}
func (u Scooter) Stop() error {
	println("stopping with" + u.Address)

	doHTTPRequest("POST", []byte(u.ID), "http://"+u.Address+"/stop")

	return nil
}

func (u Scooter) newScooter(id string, address string) Scooter {
	return Scooter{
		ID:      id,
		Address: address,
	}
}

var scooterRegistry []Scooter

func doHTTPRequest(method string, payload []byte, url string) http.Response {

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

func parseMessage(req http.Request) Message {
	var mes Message
	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&mes)
	if err != nil {
		panic(err)
	}
	if err := validate.Struct(mes); err != nil {
		panic(err)
	}
	return mes
}

func parseScooter(req http.Request) Scooter {
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

func registerScooter(s Scooter) {

	addToRegistry(s)
	fmt.Println("registred")
}

// func validateScooter(s string) Scooter {

// }

func addToRegistry(s Scooter) {
	//check if id is unique
	fmt.Println("adding to registry...")
	scooterRegistry = append(scooterRegistry, s)
}

func registerHandler(w http.ResponseWriter, req *http.Request) {
	message := parseScooter(*req)

	registerScooter(message)

}

func getEndpoints(w http.ResponseWriter, req *http.Request) {
	fmt.Println(scooterRegistry)

}

func checkScooters() {
	println("checking scooters")

}

func scooterLookup(id string) (Scooter, error) {
	fmt.Println(id)

	for _, v := range scooterRegistry {
		if v.ID == id {
			fmt.Println("found id in registry, " + id)
			return scooterRegistry[0], nil
		}
	}
	return scooterRegistry[0], fmt.Errorf("nope")

}

func startScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request

	var message = parseMessage(*req)
	// do lookup, if it within registry
	scooter, _ := scooterLookup(message.ID)

	// call that scooter and activate it
	startScooter(scooter)

}

func stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request

	var message = parseMessage(*req)
	// do lookup, if it within registry
	scooter, _ := scooterLookup(message.ID)

	// call that scooter and activate it
	stopScooter(scooter)

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
		fmt.Println("bad request")
	}
}

func getRoutes() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/start", startScooterHandler)
	http.HandleFunc("/stop", stopScooterHandler)
	http.HandleFunc("/endpoints", getEndpoints)
}

func main() {

	getRoutes()

	httpServer := http.Server{
		Addr: ":8080",
	}

	fmt.Println("serving at :8080")
	httpServer.ListenAndServe()

}
