package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/integer00/e-scooter/internal/controller/http"
	"github.com/integer00/e-scooter/internal/repository"
	webapi "github.com/integer00/e-scooter/internal/repository/webAPI"
	"github.com/integer00/e-scooter/internal/usecase"
	"github.com/integer00/e-scooter/pkg/httpserver"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	Run()
}

func Run() {

	log.Info("starting app")

	// ctx := context.Background()

	scoRegistry := repository.NewRegistry()
	scoAPP := webapi.NewScooterAPP()
	pg := repository.NewPG()
	scoUsecase := usecase.NewUseCase(scoRegistry, scoAPP, pg)
	scoController := http.NewScooterController(scoUsecase)
	//can be also like service.{component}.{component_action}

	mux := scoController.NewMux()
	handler := cors.Default().Handler(mux)

	httpServer := httpserver.New(handler)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
		httpServer.Notify()
		// case err = <-httpServer.Notify():
		// 	log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err = <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown

	// scoAPP.CallScooter()
	// scoUsecase.GetScooter("id")

}
