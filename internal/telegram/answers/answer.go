package answers

import (
	"errors"

	"github.com/whyslove/game-order-bot/internal/service"
	"github.com/whyslove/game-order-bot/internal/telegram/keyboards"
	"github.com/whyslove/game-order-bot/internal/telegram/templates"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type AnswersTelegram struct {
	bot *tgbotapi.BotAPI
}

func NewAnswers(bot *tgbotapi.BotAPI) *AnswersTelegram {
	return &AnswersTelegram{
		bot: bot,
	}
}

func (answers *AnswersTelegram) sendMessage(msg tgbotapi.MessageConfig) {
	if _, err := answers.bot.Send(msg); err != nil {
		log.Error().Msgf("error in sending message.Text: %s, err: %s", msg.Text, err.Error())
	}
}

func (answers *AnswersTelegram) AnswerError(chatID int64, err error) {
	msgText := "Произошла ошибка :("
	if err != nil {
		if errors.Is(err, service.NotAdminErr) {
			msgText = service.NotAdminErr.Error()
		}
		if errors.Is(err, service.BannedErr) {
			msgText = service.BannedErr.Error()
		}
		if errors.Is(err, service.CantDoThis) {
			msgText = service.CantDoThis.Error()
		}
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	answers.sendMessage(msg)

}

func (answers *AnswersTelegram) AnswerKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, ":)")
	msg.ReplyMarkup = keyboards.GetMenuKeyboard()
	answers.sendMessage(msg)
}

func (answers *AnswersTelegram) AnswerAdminKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, ":)")
	msg.ReplyMarkup = keyboards.GetAdminMenuKeyboard()
	answers.sendMessage(msg)
}

func (answers *AnswersTelegram) AnswerMesasge(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	answers.sendMessage(msg)
}

func (answers *AnswersTelegram) AnswerHelpMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, templates.Help)
	answers.sendMessage(msg)
}
