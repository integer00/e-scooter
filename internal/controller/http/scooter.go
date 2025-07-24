package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"github.com/integer00/e-scooter/internal/entity"
)

type Controller interface {
	NewMux() *http.ServeMux
}

var ErrMalformedJsonPayload = errors.New("malformed json or empty payload")
var ErrInvalidCookie = errors.New("cookie is invalid")

type contextMessage struct{}

var jwtkey = []byte("somesecretkey")

type ScoController struct {
	scooterService entity.ScooterService
}

func NewHTTPController(u entity.ScooterService) Controller {
	return &ScoController{
		scooterService: u,
	}
}

func (sc ScoController) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", sc.loginHandler)
	mux.HandleFunc("/registerScooter", sc.registerEndpointHandler)
	mux.HandleFunc("/scooters", sc.checkTokenMiddleware(sc.getScootersHandler))
	mux.HandleFunc("/users", sc.checkTokenMiddleware(sc.getUsersHandler))
	mux.HandleFunc("/bookscooter", sc.checkTokenMiddleware(sc.bookScooterHandler))
	mux.HandleFunc("/start", sc.checkTokenMiddleware(sc.startScooterHandler))
	mux.HandleFunc("/stop", sc.checkTokenMiddleware(sc.stopScooterHandler))
	mux.HandleFunc("/history", sc.checkTokenMiddleware(sc.rideHistoryHandler))

	return mux
}

func (sc ScoController) loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for user login")

	//get post form for user input
	// name := randomstring.HumanFriendlyEnglishString(6)
	name := new(entity.User)
	err := json.NewDecoder(req.Body).Decode(&name)
	if err != nil {
		log.Error("bad input")
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}

	user, _ := sc.scooterService.UserLogin(name.Name)
	s := generateJWT(user)

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
		// if req.Header["token"] != nil {
		// 	//verify
		// }
		// //check cookie
		// cookie, err := req.Cookie("token")
		// if err != nil {
		// 	//we have no cookie
		// 	http.Error(w, `{"status": "unauthorized"}`, http.StatusUnauthorized)
		// 	return
		// }

		// // log.Info(cookie.Value)

		// claims, valid := validateJWT(cookie.Value)
		// if !valid {
		// 	http.Error(w, ErrInvalidCookie.Error(), http.StatusBadRequest)
		// 	return
		// }

		// // should be good after, enhancing context

		// // add proper set\get for this
		// username := claims["user"].(string)

		ctx := context.WithValue(req.Context(), &contextMessage{}, "some")

		f(w, req.WithContext(ctx))
	}
}

func (sc ScoController) registerEndpointHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for registration")
	p, err := sc.parseScooter(req.Body)
	if err != nil || p == nil {
		http.Error(w, ErrMalformedJsonPayload.Error(), http.StatusBadRequest)
		return
	}

	log.Info(p)
	if ok := sc.scooterService.RegisterScooter(p); ok != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusOK)

}

func (sc ScoController) getScootersHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	//this returns scooterID+geoCoordinates like {"id":"kappa_ride","location":"coordinates"}

	log.Info("asking for endpoints")

	s := sc.scooterService.GetEndpoints()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(s)
}

func (sc ScoController) bookScooterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for booking")

	msg, err := sc.parseRequest(req.Body)
	if err != nil {
		http.Error(w, ErrMalformedJsonPayload.Error(), http.StatusBadRequest)
		return
	}
	log.Info(msg)

	if ok := sc.scooterService.BookScooter(msg.ScooterId, msg.UserId); ok != nil {
		log.Error(ok)
		http.Error(w, ok.Error(), http.StatusBadRequest)
	}
	//booked, have options to start or release

	//need common responsewriter in json

}

// might be also /api/scooter/:id/(start\stop)
func (sc ScoController) startScooterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for start")

	msg, err := sc.parseRequest(req.Body)
	if err != nil {
		http.Error(w, ErrMalformedJsonPayload.Error(), http.StatusBadRequest)
		return
	}

	if ok := sc.scooterService.StartScooter(msg.ScooterId, msg.UserId); ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
	}

}
func (sc ScoController) stopScooterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for stop")

	msg, err := sc.parseRequest(req.Body)
	if err != nil {
		http.Error(w, ErrMalformedJsonPayload.Error(), http.StatusBadRequest)
		return
	}

	if ok := sc.scooterService.StopScooter(msg.ScooterId, msg.UserId); ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
	}
}

func (sc ScoController) rideHistoryHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("asking for history")

	log.Info(req.Context())

	username := req.Context().Value(&contextMessage{}).(string)

	log.Info("username: ", username)

	sc.scooterService.RideHistory(username)

	w.WriteHeader(http.StatusOK)

}

func (sc ScoController) parseRequest(req io.Reader) (*entity.Message, error) {

	var s = new(entity.Message)
	validate := validator.New()

	err := json.NewDecoder(req).Decode(&s)
	if err != nil {
		return nil, ErrMalformedJsonPayload
	}
	if err := validate.Struct(s); err != nil {
		return nil, ErrMalformedJsonPayload
	}

	return s, nil
}

func (sc ScoController) parseScooter(req io.Reader) (*entity.Scooter, error) {

	var s = new(entity.Scooter)

	validate := validator.New()

	err := json.NewDecoder(req).Decode(&s)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	if err := validate.Struct(s); err != nil {
		log.Warn(err)
	}
	log.Info(s == nil)
	return s, nil
}

func generateJWT(name string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":  time.Now().Add(60 * time.Minute).Unix(),
		"user": name,
	})
	//todo proper claims map

	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Error("failed to sign key")
		log.Error(err)
	}
	log.Info("signing key...")

	return tokenString
}
func validateJWT(s string) (jwt.MapClaims, bool) {
	//validate token and claims
	//check expiration

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		log.Error("parse error")
		log.Error(err)
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		log.Info(claims)
		log.Info(claims["user"], claims["exp"])
	} else {
		log.Error("error with claims")
		return nil, false
	}

	return claims, true
}

func (sc ScoController) getUsersHandler(w http.ResponseWriter, req *http.Request) {
	sc.scooterService.GetUsers()
}
