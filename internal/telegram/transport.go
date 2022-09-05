package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/config"
	"github.com/whyslove/game-order-bot/internal/service"
	"github.com/whyslove/game-order-bot/internal/telegram/answers"
)

// var keyInputTeamMembers = "input_team_members"
var keyTeamID = "team_id"
var keyTeamName = "team_name"

type TelegramBot struct {
	bot         *tgbotapi.BotAPI
	answers     *answers.AnswersTelegram
	states      map[int64]string
	userStorage map[int64]map[string]string
	svc         service.ServiceI
}

func NewTelegramBot(tgToken string, svc service.ServiceI) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.ConfigValues().TelegramToken)
	if err != nil {
		return nil, err
	}
	answers := answers.NewAnswers(bot)

	states := make(map[int64]string)
	userStorage := make(map[int64]map[string]string)

	telegramBot := &TelegramBot{
		bot:         bot,
		answers:     answers,
		states:      states,
		userStorage: userStorage,
		svc:         svc,
	}

	return telegramBot, nil
}

func (tgBot *TelegramBot) StartListening(tgOffset int) {
	updateConfig := tgbotapi.NewUpdate(tgOffset) //possibly 0
	updatesCh := tgBot.bot.GetUpdatesChan(updateConfig)
	for update := range updatesCh {
		if update.CallbackQuery != nil {
			tgBot.RouteCallback(update)
		} else if update.Message != nil {
			tgBot.RouteMessage(update)
		}
		continue

	}
}

func (tgBot *TelegramBot) RefreshMatchesQueue() {
	log.Debug().Msg("Refreshing queue")
	tgBot.svc.RefreshMatches()
}

func (tgBot *TelegramBot) SetValueToUserStorage(userID int64, key string, value string) {
	if tgBot.userStorage[userID] == nil {
		tgBot.userStorage[userID] = make(map[string]string, 0)
	}

	tgBot.userStorage[userID][key] = value
}

func (tgBot *TelegramBot) GetValueFromUserStorage(userID int64, key string) string {
	if tgBot.userStorage[userID] == nil {
		return ""
	}
	return tgBot.userStorage[userID][key]
}

//TODO: доабвить конкуррентный доступ к мапе
