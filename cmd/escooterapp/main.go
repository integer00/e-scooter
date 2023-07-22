package main

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var SCOOTER_API = "http://localhost:8080/register"

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

func doRegister() error {
	jsonBody := []byte(`{"id": "kappa_ride", "address": "127.0.0.1:8081"}`)

	res := DoHTTPRequest("POST", jsonBody, SCOOTER_API)

	log.Printf("client: got response!\n")
	log.Printf("client: status code: %d\n", res.StatusCode)

	return nil
}

func getRoutes() {
	http.HandleFunc("/metrics", healthHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/start", startScooterHandler)
	http.HandleFunc("/stop", stopScooterHandler)

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

func main() {

	getRoutes()

	httpServer := http.Server{
		Addr: ":8081",
	}

	err := doRegister()
	if err != nil {
		panic("could not register itself with api")
	}

	// need to implement logic - start->go to api endpoint for registration->try until registred->send pings from time to time

	log.Println("serving at :8081")
	httpServer.ListenAndServe()

}
