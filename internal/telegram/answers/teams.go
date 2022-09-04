package answers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/whyslove/game-order-bot/internal/telegram/keyboards"
	"github.com/whyslove/game-order-bot/internal/telegram/templates"
	"github.com/whyslove/game-order-bot/internal/types"
)

func (answers *AnswersTelegram) AnswerPromptEnterTeamName(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Введите название команды!")
	answers.sendMessage(msg)
}

func (answers *AnswersTelegram) AnswerGetTeams(chatID int64, teams []types.Team) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", teams))
	answers.sendMessage(msg)

}
func (answers *AnswersTelegram) AnswerTeamCreated(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Команда создана")
	answers.sendMessage(msg)
}

func (answers *AnswersTelegram) AnswerGetMyTeams(chatID int64, teams []types.Team) {
	msg := tgbotapi.NewMessage(chatID, "*Мои команды*")
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	answers.sendMessage(msg)

	for _, team := range teams {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(templates.Team, team.Name, team.OnwerTag, team.Members))
		msg.ReplyMarkup = keyboards.GetMyTeamsKeyboard(team.Id)
		answers.sendMessage(msg)
	}
}

func (answers *AnswersTelegram) AnswerGetAllTeams(chatID int64, teams []types.Team) {
	msg := tgbotapi.NewMessage(chatID, "*Все команды*")
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	answers.sendMessage(msg)

	for _, team := range teams {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(templates.Team, team.Name, team.OnwerTag, team.Members))
		msg.ReplyMarkup = keyboards.GetMyTeamsKeyboard(team.Id)
		answers.sendMessage(msg)
	}
}
