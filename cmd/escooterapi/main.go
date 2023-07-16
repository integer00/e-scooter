package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Scooter struct {
	ID      string `json:"id" validate:"required"`
	Address string `json:"address" validate:"required"`
}

// type ScooterRegistry struct {
// 	Scooters []Scooter
// }

var scooterRegistry []Scooter

func registerHandler(w http.ResponseWriter, req *http.Request) {

	var scooter Scooter
	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&scooter)
	if err != nil {
		fmt.Println(err)
	}

	if err := validate.Struct(scooter); err != nil {
		panic(err)
	}

	addToRegistry(scooter)

	fmt.Println("registred")

}

func addToRegistry(s Scooter) {

	//check if id is unique

	scooterRegistry = append(scooterRegistry, s)

}

// func newScooter(id string, address string) Scooter {
// 	// try to register scooter
// 	// if err := json.NewDecoder(req.Body).Decode(&scooter); err !=

// 	return Scooter{
// 		ID:      "s",
// 		Address: "b",
// 	}
// }

func getEndpoints(w http.ResponseWriter, req *http.Request) {
	fmt.Println(scooterRegistry)

}

func checkScooters() {
	println("checking scooters")

}

func getRoutes() {
	http.HandleFunc("/register", registerHandler)
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
