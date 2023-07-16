package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("OK")

}

func getRoutes() {
	http.HandleFunc("/metrics", healthHandler)
	http.HandleFunc("/health", healthHandler)
}

var SCOOTER_API = "http://localhost:8080/register"

func main() {

	getRoutes()

	httpServer := http.Server{
		Addr: ":8081",
	}

	jsonBody := []byte(`{"id": "kappa_ride", "address": "127.0.0.1:8081"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, SCOOTER_API, bodyReader)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// need to implement logic - start->go to api endpoint for registration->try until registred->send pings from time to time

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	fmt.Println("serving at :8080")
	httpServer.ListenAndServe()

}
