package telegram

import (
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

var deleteTeam = regexp.MustCompile(`dteam\|.*`)

func (tgBot *TelegramBot) RouteCallback(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	chatID := update.FromChat().ID

	switch {
	case callbackData == "right_stays":
		tgBot.svc.SetMatchPlayed(false)
		matches := tgBot.svc.GetAllTodayMatches()
		tgBot.AnswerCurrentMatches(chatID, matches)
	case callbackData == "left_stays":
		tgBot.svc.SetMatchPlayed(true)
		matches := tgBot.svc.GetAllTodayMatches()
		tgBot.AnswerCurrentMatches(chatID, matches)
	case deleteTeam.MatchString(callbackData):
		res := strings.Split(callbackData, "|")
		if len(res) != 2 {
			tgBot.AnswerError(chatID)
			return
		}
		teamID, err := strconv.ParseInt(res[1], 10, 64)
		if err != nil {
			log.Error().Msgf("error converting teamID, teamID: %s, err: %s", res[1], err.Error())
			tgBot.AnswerError(chatID)
			return
		}
		err = tgBot.svc.DeleteTeam(teamID)
		if err != nil {
			log.Error().Msgf("error deleting team, teamID: %s, err: %s", res[1], err.Error())
			tgBot.AnswerError(chatID)
			return
		}
		tgBot.AnswerMesasge(chatID, "Команда удалена")
		return

	default:
		return
	}
}
