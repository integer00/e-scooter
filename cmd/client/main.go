//
//physical scooter software, communicating with api bidirectionally
//at start register itself

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/integer00/e-scooter/config"
)

var (
	QR_CODE           = "kappa_ride"
	API_ENDPOINT      = "http://localhost:8080"
	START_ENDPOINT    = API_ENDPOINT + "/start"
	STOP_ENDPOINT     = API_ENDPOINT + "/stop"
	REGISTER_ENDPOINT = API_ENDPOINT + "/registerScooter"
)

// ///controller
type IHTTPController interface {
	Run(config *config.Config)
}

type HTTPController struct {
	Scooter IScooter
}

func NewHTTPController(sc IScooter) IHTTPController {
	return &HTTPController{
		Scooter: sc,
	}
}

func (httpcontroller HTTPController) Run(config *config.Config) {

	http.HandleFunc("/metrics", healthHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/start", httpcontroller.startScooterHandler)
	http.HandleFunc("/stop", httpcontroller.stopScooterHandler)

	httpServer := http.Server{
		Addr: config.Host + ":" + config.Port,
	}

	// need to implement logic - start->go to api endpoint for registration->try until registred->send pings from time to time

	fmt.Println("serving at :8081")
	httpServer.ListenAndServe()

	// startRide()
	// time.Sleep(2 * time.Second)
	// stopRide()

}

func (httpcontroller HTTPController) startScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request
	httpcontroller.Scooter.StartScooter()

}
func (httpcontroller HTTPController) stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	// check session, parse request
	httpcontroller.Scooter.StopScooter()

}
func healthHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("OK")

}

// ///scooter
type IScooter interface {
	Bootstrap()
	RegisterScooter()
	StartScooter()
	StopScooter()
}

type Scooter struct {
	Name        string
	Id          string
	Battery     int
	Coordinates string
}

func NewScooter(name string, id string) IScooter {
	return &Scooter{
		Name: name,
		Id:   id,
	}
}

func (sc Scooter) Bootstrap() {
	jsonBody := []byte(`{"id": "kappa_ride", "address": "127.0.0.1:8081", "available": true}`)

	for {
		err := DoHTTPRequest("POST", jsonBody, REGISTER_ENDPOINT)
		fmt.Println(err)
		if err != nil {
			break
		}
		fmt.Println("could not reach bootstrap server, retrying")
		time.Sleep(5 * time.Second)
	}

	fmt.Printf("client: registred with server")
}
func (sc Scooter) RegisterScooter() {

}
func (sc Scooter) StartScooter() {
	fmt.Println("scooter is started!")
}
func (sc Scooter) StopScooter() {
	fmt.Println("scooter is stopped!")
}

func DoHTTPRequest(method string, payload []byte, url string) *http.Response {

	bodyReader := bytes.NewReader(payload)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		println("request failed")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return res
}

func main() {

	scooter := NewScooter("kappa_ride", "553")

	controller := NewHTTPController(scooter)

	scooter.Bootstrap()

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8081")

	config := config.NewConfig()

	controller.Run(config)

}

// func parseQRCode() string {
// 	//here's goes camera parsing real qr-code
// 	return QR_CODE
// }

// func startRide() {
// 	jsonBody := []byte(`{"id": "kappa_ride"}`)
// 	var url = START_ENDPOINT

// 	response := DoHTTPRequest("POST", jsonBody, url)

// 	r, _ := ioutil.ReadAll(response.Body)

// 	println(response.StatusCode, string(r))
// }

// func stopRide() {
// 	jsonBody := []byte(`{"id": "kappa_ride"}`)
// 	var url = STOP_ENDPOINT

// 	response := DoHTTPRequest("POST", jsonBody, url)

// 	r, _ := ioutil.ReadAll(response.Body)

// 	println(response.StatusCode, string(r))
// }
