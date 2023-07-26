package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator"
	"github.com/integer00/e-scooter/internal/entity"
)

type ScoController struct {
	scooterUseCase entity.UseCase
}

func NewScooterController(u entity.UseCase) entity.Controller {
	return &ScoController{
		scooterUseCase: u,
	}
}

func (sc ScoController) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", sc.loginHandler)
	mux.HandleFunc("/registerScooter", sc.registerEndpointHandler)
	mux.HandleFunc("/scooters", sc.checkTokenMiddleware(sc.getScootersHandler))
	mux.HandleFunc("/start", sc.startScooterHandler)
	mux.HandleFunc("/stop", sc.stopScooterHandler)
	return mux
}

func (sc ScoController) loginHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for user login")

	//get post form for user input
	// name := randomstring.HumanFriendlyEnglishString(6)

	s, _ := sc.scooterUseCase.UserLogin("alice")

	cookie := http.Cookie{
		Name:    "token",
		Domain:  "localhost",
		Path:    "/",
		Expires: time.Now().Add(10 * time.Minute),
		Value:   s,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, req, "http://localhost:8080/scooters", http.StatusMovedPermanently)

	log.Info(s)

}

func (sc ScoController) checkTokenMiddleware(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Header["token"] != nil {
			//verify
		}
		//check cookie
		cookie, err := req.Cookie("token")
		if err != nil {
			//we have no cookie
			http.Error(w, `{"status": "unauthorized"}`, http.StatusUnauthorized)
			return

		}
		// log.Info(cookie.Value)

		sc.scooterUseCase.ValidateJWT(cookie.Value)

		f(w, req)
	}

}

func (sc ScoController) registerEndpointHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for registration")
	p := parseRequest(*req)

	sc.scooterUseCase.RegisterScooter(p)

	w.WriteHeader(http.StatusOK)

}

func (sc ScoController) getScootersHandler(w http.ResponseWriter, req *http.Request) {
	//this returns scooterID+geoCoordinates like {"id":"kappa_ride","location":"coordinates"}

	log.Info("asking for endpoints")

	s := sc.scooterUseCase.GetEndpoints()
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(s)

}

// might be also /api/scooter/:id/(start\stop)
func (sc ScoController) startScooterHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for start")

	s := parseRequest(*req)

	ctx := context.WithValue(context.Background(), "scooter", s.Id)

	err := sc.scooterUseCase.StartScooter(ctx)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

}
func (sc ScoController) stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("asking for stop")

	s := parseRequest(*req)

	ctx := context.WithValue(context.Background(), "scooter", s.Id)

	err := sc.scooterUseCase.StopScooter(ctx)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func parseRequest(req http.Request) entity.Scooter {
	var s = entity.Scooter{}

	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&s)
	if err != nil {
		log.Warn(err)

	}
	if err := validate.Struct(s); err != nil {
		log.Warn(err)
	}
	return s
}
