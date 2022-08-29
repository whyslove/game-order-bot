package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	stateInputCommandName = "input_command_name"
	stateEmpty            = ""
)

func (tgBot *TelegramBot) RouteMessage(update tgbotapi.Update) {
	userID := update.Message.From.ID
	state := tgBot.states[userID]

	// Handle commands
	if update.Message.IsCommand() { // ignore any non-command Messages
		switch update.Message.Command() {
		case "keyboard":
			tgBot.AnswerKeyboard(update.Message.Chat.ID)
			return
		case "start":
			tgBot.AnswerMesasge(update.Message.Chat.ID, "Привет! Напиши /keyboard чтобы появилась клавиатура. Напиши /help для помощи")
		case "stop":
			tgBot.states[userID] = ""
			tgBot.AnswerMesasge(update.Message.Chat.ID, "Состояние сброшено")
		case "help":
			tgBot.AnswerHelpMessage(update.Message.Chat.ID)
		default:
			tgBot.AnswerMesasge(update.Message.Chat.ID, "Команда не распознана")
		}
	}

	// Handle scpecific states
	switch state {
	case stateInputCommandName:
		err := tgBot.svc.CreateTeam(update.Message.Text, userID, update.Message.From.UserName)
		tgBot.states[userID] = stateEmpty
		if err != nil {
			log.Error().Msgf("error in stateInputCommandName, err: %s", err.Error())
			tgBot.AnswerError(update.Message.Chat.ID)
			return
		}
		tgBot.AnswerTeamCreated(update.Message.Chat.ID)
		return
	}

	// Handle other
	switch update.Message.Text {
	case "Расписание":
		matches := tgBot.svc.GetAllTodayMatches()
		tgBot.AnswerCurrentMatches(update.Message.Chat.ID, matches)
	case "Новая команда":
		tgBot.states[userID] = stateInputCommandName
		tgBot.AnswerPromptEnterTeamName(update.Message.Chat.ID)
	case "Все команды":
		teams, err := tgBot.svc.GetTodayTeams()
		if err != nil {
			log.Error().Msgf("error in Все команды, err: %s", err.Error())
			tgBot.AnswerError(update.Message.Chat.ID)
			return
		}
		tgBot.AnswerGetTeams(update.Message.Chat.ID, teams)
	case "Мои команды":
		teams, err := tgBot.svc.GetMyTeams(userID)
		if err != nil {
			log.Error().Msgf("error getting my teams err: %s", err.Error())
			tgBot.AnswerError(update.Message.Chat.ID)
			return
		}
		tgBot.AnswerGetMyTeams(update.Message.Chat.ID, teams)
	case "Удалить все":
		err := tgBot.svc.DeleteAllInformationToday()
		if err != nil {
			log.Error().Msgf("error deleting all err: %s", err.Error())
			tgBot.AnswerError(update.Message.Chat.ID)
			return
		}
		tgBot.AnswerMesasge(update.Message.Chat.ID, "Все удалено")
	}
}
