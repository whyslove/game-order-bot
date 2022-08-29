package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/config"
	"github.com/whyslove/game-order-bot/internal/service"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	states map[int64]string
	svc    service.ServiceI
}

func NewTelegramBot(tgToken string, svc service.ServiceI) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.ConfigValues().TelegramToken)
	if err != nil {
		return nil, err
	}

	states := make(map[int64]string)

	telegramBot := &TelegramBot{
		bot:    bot,
		states: states,
		svc:    svc,
	}

	return telegramBot, nil
}

func (tgBot *TelegramBot) StartListening(tgOffset int) {
	updateConfig := tgbotapi.NewUpdate(tgOffset) //possibly 0
	updatesCh := tgBot.bot.GetUpdatesChan(updateConfig)
	for update := range updatesCh {
		if update.CallbackQuery != nil {
			log.Debug().Msg("here")
			tgBot.RouteCallback(update)
		} else if update.Message != nil {
			tgBot.RouteMessage(update)
		}
		continue

	}
}

//TODO: доабвить конкуррентный доступ к мапе
