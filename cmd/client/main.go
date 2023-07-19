package main

import (
	"time"

	"github.com/integer00/e-scooter/internal/repo"
)

var (
	QR_CODE        = "kappa_ride"
	API_ENDPOINT   = "http://localhost:8080"
	START_ENDPOINT = API_ENDPOINT + "/start"
	STOP_ENDPOINT  = API_ENDPOINT + "/stop"
)

func parseQRCode() string {
	//here's goes camera parsing real qr-code
	return QR_CODE
}

func startRide() {
	jsonBody := []byte(`{"id": "kappa_ride"}`)
	var url = START_ENDPOINT

	response := repo.DoHTTPRequest("POST", jsonBody, url)

	println(response.StatusCode)
}

func stopRide() {
	jsonBody := []byte(`{"id": "kappa_ride"}`)
	var url = STOP_ENDPOINT

	response := repo.DoHTTPRequest("POST", jsonBody, url)

	println(response.StatusCode)
}

func main() {

	startRide()
	time.Sleep(10 * time.Second)
	stopRide()

}
