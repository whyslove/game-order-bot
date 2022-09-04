package answers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/whyslove/game-order-bot/internal/telegram/keyboards"
	"github.com/whyslove/game-order-bot/internal/telegram/templates"
	"github.com/whyslove/game-order-bot/internal/types"
)

func (answers *AnswersTelegram) AnswerCurrentMatches(chatID int64, matches []types.MatchQueue) {
	msg := tgbotapi.NewMessage(chatID, templates.MatchHeader)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	answers.sendMessage(msg)

	for _, match := range matches {
		score := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, match.Score)
		t1 := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, match.Team1)
		t2 := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, match.Team2)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(templates.Match, t1, t2, score))
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if match.Current {
			msg.ReplyMarkup = keyboards.GetCurrentMatchKeyboard()
		}
		answers.sendMessage(msg)
	}
}
