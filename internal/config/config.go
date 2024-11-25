package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	BindAddress string

	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	APIUrl  string
	APIPort string
}

func NewConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file")
	}

	log.Debug("environment variables loaded")

	config := Config{
		BindAddress:      os.Getenv("BIND_ADDRESS"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		APIUrl:           os.Getenv("API_URL"),
		APIPort:          os.Getenv("API_PORT"),
	}

	return config
}
