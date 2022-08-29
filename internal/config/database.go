package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type databaseStruct struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
	User     string `envconfig:"DB_USER" required:"true"`
	DbName   string `envconfig:"DB_NAME" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	SslMode  string `envconfig:"DB_SSL_MODE" required:"true"`
}

var database *databaseStruct

func DatabaseValues() *databaseStruct {
	if database != nil {
		return database
	}

	database := &databaseStruct{}
	godotenv.Load(".env")

	if err := envconfig.Process("", database); err != nil {
		log.Fatal().Msgf("error in processing database config values, err: %s", err.Error())
	}
	return database
}
