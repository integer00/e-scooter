package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
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

	response := DoHTTPRequest("POST", jsonBody, url)

	r, _ := ioutil.ReadAll(response.Body)

	println(response.StatusCode, string(r))
}

func stopRide() {
	jsonBody := []byte(`{"id": "kappa_ride"}`)
	var url = STOP_ENDPOINT

	response := DoHTTPRequest("POST", jsonBody, url)

	r, _ := ioutil.ReadAll(response.Body)

	println(response.StatusCode, string(r))
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

	startRide()
	time.Sleep(2 * time.Second)
	stopRide()

}
