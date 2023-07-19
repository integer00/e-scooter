package main

import (
	"net/http"

	"github.com/integer00/e-scooter/internal/repo"
	log "github.com/sirupsen/logrus"
)

func registerHandler(w http.ResponseWriter, req *http.Request) {
	message := repo.ParseScooter(*req)

	repo.RegisterScooter(message)
}

func getEndpoints(w http.ResponseWriter, req *http.Request) {
	// var response = scooterRegistry.(string)

	// log.Fprintf(w, scooterRegistry)
	log.Println("endpoints")
}

func startScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request

	var message = repo.ParseMessage(*req)
	// do lookup, if it within registry

	//charge for a fee
	err := repo.PgHandler("userid", "startRide")
	if err != nil {
		log.Println("error with PG")
		panic(err)
	}

	repo.ScooterHandler(message, "startAction")

}

func stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request

	var message = repo.ParseMessage(*req)

	//charge for a Fair, returning deposit
	err := repo.PgHandler("userid", "endRide")
	if err != nil {
		log.Println("error with PG")
		panic(err)
	}

	// call that scooter and activate it
	repo.ScooterHandler(message, "stopAction")

}

func getRoutes() {
	http.HandleFunc("/register", registerHandler)
	/*
		handle register requests, send them to register usecase
	*/
	http.HandleFunc("/start", startScooterHandler)
	http.HandleFunc("/stop", stopScooterHandler)
	http.HandleFunc("/endpoints", getEndpoints)
}

func main() {

	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{
	// 	FullTimestamp: true,
	// })

	getRoutes()

	httpServer := http.Server{
		Addr: ":8080",
	}

	log.Println("starting at :8080")
	httpServer.ListenAndServe()

}
