package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type configStruct struct {
	TelegramToken string `envconfig:"TG_TOKEN" required:"true"`
}

var config *configStruct

func ConfigValues() configStruct {

	if config != nil {
		return *config
	}

	config := &configStruct{}
	godotenv.Load(".env")

	if err := envconfig.Process("", config); err != nil {
		log.Fatal().Msgf("error in processing config, err: %s", err.Error())
	}

	return *config
}
