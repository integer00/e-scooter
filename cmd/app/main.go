package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/integer00/e-scooter/config"
	"github.com/integer00/e-scooter/internal/controller/http"
	"github.com/integer00/e-scooter/internal/registry"
	"github.com/integer00/e-scooter/internal/repository"
	"github.com/integer00/e-scooter/internal/service"
	"github.com/integer00/e-scooter/pkg/httpserver"
	"github.com/rs/cors"
)

func main() {
	Run()
	Shutdown()
}

func Run() {

	println("starting app")

	//should be in ENV within your container, setting here for convenience
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("PG_POOL_MAX", "10")
	os.Setenv("PG_URL", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")

	config := config.NewConfig()

	scooterRegistry := registry.NewRegistry()
	userRegistry := registry.NewUserRegistry()
	scooterRepository := repository.NewScooterRepository()
	paymentGate := repository.NewPG()

	// log.Info("getting postgres")
	// pg, err := postgres.New(config.PG_URL)
	// if err != nil {
	// 	log.Fatal("failed to initialize postgresql connection")
	// }
	// defer pg.Close()

	// pgRepo := repository.NewPostgresRepo(pg) //new database, use interface , and in CLI argument when start use USE_PG=1, or fallover to internal inmem

	scooterService := service.NewService(scooterRegistry, scooterRepository, paymentGate, userRegistry, nil)

	scooterHTTPController := http.NewHTTPController(scooterService)

	mux := scooterHTTPController.NewMux()

	//cors needs to be adjusted in production environment
	handler := cors.Default().Handler(mux)

	httpServer := httpserver.New(handler, config.Host+":"+config.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	//fix loop
	select {
	case s := <-interrupt:
		print("app - Run - signal: " + s.String())
		httpServer.Notify()
		// case err = <-httpServer.Notify():
		// 	log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err = <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

}

func Shutdown() {
	println("Cleanup after exit")
}
