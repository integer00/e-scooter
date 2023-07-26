package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/integer00/e-scooter/config"
	log "github.com/sirupsen/logrus"
)

var SCOOTER_API = "http://localhost:8080/registerScooter"

func startScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request
	startScooter()

}
func stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request
	stopScooter()

}

func startScooter() {
	log.Println("scooter is started!")
}
func stopScooter() {
	log.Println("scooter is stopped!")
}

func healthHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("OK")

}

func bootstrap() {
	jsonBody := []byte(`{"id": "kappa_ride", "address": "127.0.0.1:8081"}`)

	for {
		err := DoHTTPRequest("POST", jsonBody, SCOOTER_API)
		if err == nil {
			break
		}
		log.Warn("could not reach bootstrap server, retrying")
		time.Sleep(5 * time.Second)
	}

	log.Printf("client: registred with server")

}

func getRoutes() {
	http.HandleFunc("/metrics", healthHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/start", startScooterHandler)
	http.HandleFunc("/stop", stopScooterHandler)

}

func DoHTTPRequest(method string, payload []byte, url string) error {

	bodyReader := bytes.NewReader(payload)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return errors.New("could not register")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("could not register")
	}
	log.Info(res.StatusCode)

	return nil
}

func main() {

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8081")

	config := config.NewConfig()

	getRoutes()

	httpServer := http.Server{
		Addr: config.Host + ":" + config.Port,
	}

	bootstrap()

	// need to implement logic - start->go to api endpoint for registration->try until registred->send pings from time to time

	log.Println("serving at :8081")
	httpServer.ListenAndServe()

}
