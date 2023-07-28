package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/integer00/e-scooter/config"
	"github.com/integer00/e-scooter/internal/controller/http"
	"github.com/integer00/e-scooter/internal/repository"
	webapi "github.com/integer00/e-scooter/internal/repository/webAPI"
	"github.com/integer00/e-scooter/internal/usecase"
	"github.com/integer00/e-scooter/pkg/httpserver"
	"github.com/integer00/e-scooter/pkg/postgres"
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
	os.Setenv("PG_POOL_MAX", "10")
	os.Setenv("PG_URL", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")

	config := config.NewConfig()

	scoRegistry := repository.NewRegistry()
	userRegistry := repository.NewUserRegistry()
	scoAPP := webapi.NewScooterAPP()
	pgate := repository.NewPG()

	log.Info("getting postgres")
	pg, err := postgres.New(config.PG_URL)
	if err != nil {
		log.Fatal("failed to initialize postgresql connection")
	}
	defer pg.Close()

	pgRepo := repository.NewPostgresRepo(pg)

	var id int
	var name string
	log.Info(id, name)

	er := pg.Pool.QueryRow(context.Background(), "select * from users").Scan(&id, &name)
	if er != nil {
		log.Info(er)
	}

	log.Info(id, name)

	scoUsecase := usecase.NewUseCase(scoRegistry, scoAPP, pgate, userRegistry, pgRepo)
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
