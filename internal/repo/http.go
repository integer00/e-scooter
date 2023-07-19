package repo

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/integer00/e-scooter/internal/models"
)

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

func ParseMessage(req http.Request) models.Message {
	var mes models.Message
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
