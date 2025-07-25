package main

import (
	"fmt"
	"os"

	"github.com/integer00/e-scooter/config"
	"github.com/integer00/e-scooter/internal/controller/http"
	"github.com/integer00/e-scooter/internal/registry"
	"github.com/integer00/e-scooter/internal/repository"
	"github.com/integer00/e-scooter/internal/service"
)

func main() {
	Run()
	Shutdown()
}

func Run() {

	fmt.Println("starting app")

	//should be in ENV within your container, setting here for convenience
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("PG_POOL_MAX", "10")
	os.Setenv("PG_URL", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")

	config := config.NewConfig()

	scooterRegistry := registry.NewRegistry()
	scooterRepository := repository.NewScooterRepository()
	paymentGate := repository.NewPaymentGate()

	// log.Info("getting postgres")
	// pg, err := postgres.New(config.PG_URL)
	// if err != nil {
	// 	log.Fatal("failed to initialize postgresql connection")
	// }
	// defer pg.Close()

	// pgRepo := repository.NewPostgresRepo(pg) //new database, use interface , and in CLI argument when start use USE_PG=1, or fallover to internal inmem

	scooterService := service.NewService(scooterRegistry, scooterRepository, paymentGate, nil)

	scooterHTTPController := http.NewHTTPController(scooterService)

	scooterHTTPController.Run(config)
	//GRPCController.Run(config)

}

func Shutdown() {
	fmt.Println("Cleanup after exit")
}
