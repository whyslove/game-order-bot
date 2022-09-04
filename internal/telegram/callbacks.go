package telegram

import (
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

var deleteTeam = regexp.MustCompile(`dteam\|.*`)
var changeTeamMembers = regexp.MustCompile(`chteammemebers\|.*`)

func (tgBot *TelegramBot) RouteCallback(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	chatID := update.FromChat().ID
	userID := update.SentFrom().ID

	switch {
	case callbackData == "right_stays":
		err := tgBot.svc.SetMatchPlayed(userID, false)
		if err != nil {
			log.Error().Msgf("userID %d, err: %s", userID, err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		matches := tgBot.svc.GetAllTodayMatches()
		tgBot.answers.AnswerCurrentMatches(chatID, matches)
	case callbackData == "left_stays":
		err := tgBot.svc.SetMatchPlayed(userID, true)
		if err != nil {
			log.Error().Msgf("userID %d, err: %s", userID, err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		matches := tgBot.svc.GetAllTodayMatches()
		tgBot.answers.AnswerCurrentMatches(chatID, matches)
	case deleteTeam.MatchString(callbackData):
		res := strings.Split(callbackData, "|")
		if len(res) != 2 {
			tgBot.answers.AnswerError(chatID, nil)
			return
		}
		teamID, err := strconv.ParseInt(res[1], 10, 64)
		if err != nil {
			log.Error().Msgf("error converting teamID, teamID: %s, err: %s", res[1], err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		err = tgBot.svc.DeleteTeam(userID, teamID)
		if err != nil {
			log.Error().Msgf("error deleting team, teamID: %s, err: %s", res[1], err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}
		tgBot.answers.AnswerMesasge(chatID, "Команда удалена")
		return
	case changeTeamMembers.MatchString(callbackData):
		res := strings.Split(callbackData, "|")
		if len(res) != 2 {
			tgBot.answers.AnswerError(chatID, nil)
			return
		}
		_, err := strconv.ParseInt(res[1], 10, 64)
		if err != nil {
			log.Error().Msgf("error converting teamID, teamID: %s, err: %s", res[1], err.Error())
			tgBot.answers.AnswerError(chatID, err)
			return
		}

		tgBot.states[userID] = stateUpdatingTeamMembers
		tgBot.SetValueToUserStorage(userID, keyTeamID, res[1]) //res[1] is teamID encoded in callback data

		tgBot.answers.AnswerMesasge(chatID, "Введите новый состав")

	default:
		return
	}
}
