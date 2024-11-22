package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	BindAddress string `env:"BIND_ADDRESS" env-default:":8080"`

	PostgresHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" env-default:"postgres"`
	PostgresUser     string `env:"POSTGRES_USER" env-default:"admin"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-default:"admin"`

	APIURL string `env:"API_URL" env-default:"http://localhost:8080"`
}

func NewConfig() Config {
	config := Config{}
	if err := cleanenv.ReadEnv(&config); err != nil {
		panic(fmt.Errorf("error reading config: %w", err))
	}

	log.Infof("config loaded %+v", config)

	return config
}
