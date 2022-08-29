package telegram

import (
	"fmt"

	"github.com/whyslove/game-order-bot/internal/telegram/keyboards"
	"github.com/whyslove/game-order-bot/internal/telegram/templates"
	"github.com/whyslove/game-order-bot/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (tgBot *TelegramBot) sendMessage(msg tgbotapi.MessageConfig) {
	if _, err := tgBot.bot.Send(msg); err != nil {
		log.Error().Msgf("error in sending message.Text: %s, err: %s", msg.Text, err.Error())
	}
}

func (tgBot *TelegramBot) AnswerPromptEnterTeamName(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Введите название команды!")
	tgBot.sendMessage(msg)
}

func (tgBot *TelegramBot) AnswerGetTeams(chatID int64, teams []types.Team) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", teams))
	tgBot.sendMessage(msg)

}

func (tgBot *TelegramBot) AnswerError(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Произошла ошибка :(")
	tgBot.sendMessage(msg)

}

func (tgBot *TelegramBot) AnswerTeamCreated(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Команда создана")
	tgBot.sendMessage(msg)
}

func (tgBot *TelegramBot) AnswerCurrentMatches(chatID int64, matches []types.MatchQueue) {
	msg := tgbotapi.NewMessage(chatID, "*Расписание:*")
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	tgBot.sendMessage(msg)

	for _, match := range matches {
		score := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, match.Score)
		t1 := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, match.Team1)
		t2 := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, match.Team2)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(templates.Match, t1, t2, score))
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if match.Current {
			msg.ReplyMarkup = keyboards.GetCurrentMatchKeyboard()
		}
		tgBot.sendMessage(msg)
	}
}

func (tgBot *TelegramBot) AnswerKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, ":)")
	msg.ReplyMarkup = keyboards.GetMenuKeyboard()
	tgBot.sendMessage(msg)
}

func (tgBot *TelegramBot) AnswerGetMyTeams(chatID int64, teams []types.Team) {
	msg := tgbotapi.NewMessage(chatID, "*Мои команды*")
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	tgBot.sendMessage(msg)

	for _, team := range teams {
		msg := tgbotapi.NewMessage(chatID, team.Name)
		msg.ReplyMarkup = keyboards.GetMyTeamsKeyboard(team.Id)
		tgBot.sendMessage(msg)
	}
}

func (tgBot *TelegramBot) AnswerMesasge(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	tgBot.sendMessage(msg)
}
