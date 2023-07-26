package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/integer00/e-scooter/config"
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

	//should be in ENV within your container, setting here for convenience
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")

	// ctx := context.Background()

	config := config.NewConfig()

	scoRegistry := repository.NewRegistry()
	userRegistry := repository.NewUserRegistry()
	scoAPP := webapi.NewScooterAPP()
	pg := repository.NewPG()

	scoUsecase := usecase.NewUseCase(scoRegistry, scoAPP, pg, userRegistry)
	scoController := http.NewScooterController(scoUsecase)
	//can be also like service.{component}.{component_action}

	mux := scoController.NewMux()

	//cors needs to be adjusted in production environment
	handler := cors.Default().Handler(mux)

	httpServer := httpserver.New(handler, config.Host+":"+config.Port)

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
