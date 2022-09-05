package main

import (
	"time"

	"github.com/robfig/cron"
	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/config"
	"github.com/whyslove/game-order-bot/internal/service"
	"github.com/whyslove/game-order-bot/internal/storage"
	"github.com/whyslove/game-order-bot/internal/telegram"
)

var tgOffset = 0

func main() {
	db, err := storage.NewPostgresDb(config.DatabaseValues().Host, config.DatabaseValues().Port, config.DatabaseValues().User,
		config.DatabaseValues().DbName, config.DatabaseValues().Password, config.DatabaseValues().SslMode)
	if err != nil {
		log.Fatal().Msgf("error initializing database, err: %s", err.Error())
	}

	svc := service.NewService(db)

	tgBot, err := telegram.NewTelegramBot(config.ConfigValues().TelegramToken, svc)
	if err != nil {
		log.Fatal().Msgf("error in creating tgBot err: %s", err.Error())
	}

	customLocation, _ := time.LoadLocation("Europe/Moscow")
	log.Info().Msgf("location: %v", customLocation.String())
	cronHandler := cron.NewWithLocation(customLocation)
	cronHandler.AddFunc("* * * * *", func() {
		tgBot.RefreshMatchesQueue()
	})
	cronHandler.Start()
	tgBot.StartListening(tgOffset)
}
