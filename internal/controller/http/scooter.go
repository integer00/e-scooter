package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"

	"github.com/integer00/e-scooter/internal/entity"
)

type ScoController struct {
	scooterUseCase entity.UseCase
	context        context.Context
}

type contextHandler struct {
	context context.Context
	handler http.Handler
}

func NewScooterController(u entity.UseCase) entity.Controller {
	return &ScoController{
		scooterUseCase: u,
		context:        context.Background(),
	}
}

func (sc ScoController) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", sc.loginHandler)
	mux.HandleFunc("/registerScooter", sc.registerEndpointHandler)
	mux.HandleFunc("/scooters", sc.checkTokenMiddleware(sc.getScootersHandler))
	mux.HandleFunc("/bookscooter", sc.checkTokenMiddleware(sc.bookScooterHandler))
	mux.HandleFunc("/start", sc.checkTokenMiddleware(sc.startScooterHandler))
	mux.HandleFunc("/stop", sc.checkTokenMiddleware(sc.stopScooterHandler))
	return mux
}

func (sc ScoController) loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for user login")

	//get post form for user input
	// name := randomstring.HumanFriendlyEnglishString(6)

	s, _ := sc.scooterUseCase.UserLogin("alice")

	cookie := http.Cookie{
		Name:    "token",
		Domain:  "localhost",
		Path:    "/",
		Expires: time.Now().Add(60 * time.Minute),
		Value:   s,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, req, "http://localhost:8080/scooters", http.StatusMovedPermanently)

	log.Info(s)

}

func (sc *ScoController) checkTokenMiddleware(f http.HandlerFunc) http.HandlerFunc {

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

		if !sc.scooterUseCase.ValidateJWT(cookie.Value) {
			http.Error(w, "cookie is invalid", http.StatusBadRequest)
			return
		}

		f(w, req)
	}

}

func (sc ScoController) registerEndpointHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for registration")
	p := parseScooter(*req)

	sc.scooterUseCase.RegisterScooter(p)

	w.WriteHeader(http.StatusOK)

}

func (sc ScoController) bookScooterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for booking")

	msg := parseRequest(*req)
	if msg == nil {
		http.Error(w, "malformed json or empty payload", http.StatusBadRequest)
		return
	}
	log.Info(msg)

	err := sc.scooterUseCase.BookScooter("kappa_ride", "alice")
	if err != nil {
		log.Error(err)
	}
	//booked, have options to start or release

}

func (sc ScoController) getScootersHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	//this returns scooterID+geoCoordinates like {"id":"kappa_ride","location":"coordinates"}

	log.Info("asking for endpoints")

	s := sc.scooterUseCase.GetEndpoints()
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(s)

}

// might be also /api/scooter/:id/(start\stop)
func (sc ScoController) startScooterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for start")

	msg := parseRequest(*req)
	if msg == nil {
		http.Error(w, "malformed json or empty payload", http.StatusBadRequest)
		return
	}

	err := sc.scooterUseCase.StartScooter(msg.Scooterid)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

}
func (sc ScoController) stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for stop")
	//getting ride

	// ride :=

	msg := parseRequest(*req)
	if msg == nil {
		http.Error(w, "malformed json or empty payload", http.StatusBadRequest)
		return

	}

	err := sc.scooterUseCase.StopScooter(msg.Scooterid)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func parseRequest(req http.Request) *entity.Message {

	var s = new(entity.Message)

	// validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&s)
	if err != nil {
		return nil
	}
	// if err := validate.Struct(s); err != nil {
	// 	log.Warn(err)
	// }
	return s
}

func parseScooter(req http.Request) *entity.Scooter {

	var s = new(entity.Scooter)

	validate := validator.New()

	err := json.NewDecoder(req.Body).Decode(&s)
	if err != nil {
		return nil
	}
	if err := validate.Struct(s); err != nil {
		log.Warn(err)
	}
	return s
}
