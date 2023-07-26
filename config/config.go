package config

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host string `env:"HOST" env-required:"true"`
	Port string `env:"PORT" env-required:"true"`
}

func NewConfig() *Config {
	config := &Config{}
	log.Info(os.Getenv("HOST"))
	log.Info(os.Getenv("PORT"))

	err := cleanenv.ReadEnv(config)
	if err != nil {
		log.Error(err)
		log.Fatal("failed to read config")
	}

	return config
}
