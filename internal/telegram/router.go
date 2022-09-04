package telegram

import (
	"errors"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/service"
	"github.com/whyslove/game-order-bot/internal/types"
)

const (
	stateEmpty                      = ""
	stateInputCommandName           = "input_team_name"
	stateInputTokenForCreatingAdmin = "input_token_for_become_admin"
	stateUpdatingTeamMembers        = "update_team_members"
	stateInputTeamMembers           = "input_team_members"
)

func (tgBot *TelegramBot) RouteMessage(update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	state := tgBot.states[userID]

	// Handle commands
	if update.Message.IsCommand() { // ignore any non-command Messages
		switch update.Message.Command() {
		case "keyboard":
			verdict, err := tgBot.svc.CheckIsAdmin(userID)
			if err != nil {
				log.Error().Msgf("error checking is admin, err: %s", err.Error())
				tgBot.answers.AnswerKeyboard(chatID)
			} else if !verdict {
				tgBot.answers.AnswerKeyboard(chatID)
			} else {
				tgBot.answers.AnswerAdminKeyboard(chatID)
			}
			return
		case "start":
			tgBot.answers.AnswerMesasge(chatID, "Привет! Напиши /keyboard чтобы появилась клавиатура. Напиши /help для помощи")
		case "stop":
			tgBot.states[userID] = ""
			tgBot.answers.AnswerMesasge(chatID, "Состояние сброшено")
		case "help":
			tgBot.answers.AnswerHelpMessage(chatID)
		default:
			tgBot.answers.AnswerMesasge(chatID, "Команда не распознана")
		}
		return
	}

	// Handle states
	switch state {
	case stateInputCommandName:
		tgBot.states[userID] = stateInputTeamMembers
		tgBot.SetValueToUserStorage(userID, keyTeamName, update.Message.Text)
		tgBot.answers.AnswerMesasge(chatID, "Введите участников команды")
		return
	case stateInputTokenForCreatingAdmin:
		tgBot.states[userID] = stateEmpty
		user := types.User{
			TgID:     userID,
			Name:     update.Message.From.UserName,
			IsBanned: false,
			IsAdmin:  true,
		}
		err := tgBot.svc.SaveUser(userID, update.Message.Text, user)
		if err != nil {
			log.Error().Msgf("error in stateInputTokenForCreatingAdmin, err: %s", err.Error())
			if errors.Is(err, service.ErrbadToken) {
				tgBot.answers.AnswerMesasge(chatID, "Неправильный / уже использованный токен")
				return
			}
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerMesasge(chatID, "Вы получили права админа :)")
		tgBot.answers.AnswerAdminKeyboard(chatID)
		return

	case stateInputTeamMembers:
		tgBot.states[userID] = stateEmpty
		defer func() {
			tgBot.SetValueToUserStorage(userID, keyTeamName, "")
		}()

		teamName := tgBot.GetValueFromUserStorage(userID, keyTeamName)
		err := tgBot.svc.CreateTeam(userID, teamName, userID, update.Message.From.UserName, update.Message.Text)
		if err != nil {
			log.Error().Msgf("error in stateInputCommandName, err: %s", err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerTeamCreated(chatID)
		return
	case stateUpdatingTeamMembers:
		tgBot.states[userID] = stateEmpty
		defer func() {
			tgBot.SetValueToUserStorage(userID, keyTeamID, "")
		}()

		strTeamID := tgBot.GetValueFromUserStorage(userID, keyTeamID)
		teamID, err := strconv.ParseInt(strTeamID, 10, 64)
		if err != nil {
			log.Error().Msgf("error converting teamID, teamID: %s, err: %s", strTeamID, err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}

		err = tgBot.svc.UpdateTeamMembers(userID, teamID, update.Message.Text)
		if err != nil {
			log.Error().Msgf("error updating team member teamID, teamID: %s, err: %s", strTeamID, err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerMesasge(chatID, "Состав обновлен")
		return
	}

	// Handle other
	switch update.Message.Text {
	case "Расписание":
		matches := tgBot.svc.GetAllTodayMatches()
		tgBot.answers.AnswerCurrentMatches(chatID, matches)
	case "Новая команда":
		tgBot.states[userID] = stateInputCommandName
		tgBot.answers.AnswerPromptEnterTeamName(chatID)
	case "Все команды":
		teams, err := tgBot.svc.GetTodayTeams(userID)
		if err != nil {
			log.Error().Msgf("error in Все команды, err: %s", err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerGetAllTeams(chatID, teams)
	case "Мои команды":
		teams, err := tgBot.svc.GetMyTeams(userID)
		if err != nil {
			log.Error().Msgf("error getting my teams err: %s", err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerGetMyTeams(chatID, teams)
	case "Удалить все":
		err := tgBot.svc.DeleteAllInformationToday(userID)
		if err != nil {
			log.Error().Msgf("error deleting all err: %s", err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerMesasge(chatID, "Все удалено")
	case "Ввести токен":
		tgBot.states[userID] = stateInputTokenForCreatingAdmin
		tgBot.answers.AnswerMesasge(chatID, "Введите токен")
	}

}
