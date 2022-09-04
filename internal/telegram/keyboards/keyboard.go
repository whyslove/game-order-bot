package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Расписание"),
			tgbotapi.NewKeyboardButton("Мои команды"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Новая команда"),
		),
	)
}

func GetAdminMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Расписание"),
			tgbotapi.NewKeyboardButton("Все команды"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Новая команда"),
			tgbotapi.NewKeyboardButton("Удалить все")),
	)
}

func GetCurrentMatchKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Правая команда остается", "right_stays"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Левая команда остается", "left_stays"),
		),
	)
}

func GetMyTeamsKeyboard(teamID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить команду", fmt.Sprintf("dteam|%d", teamID)),
			tgbotapi.NewInlineKeyboardButtonData("Изменить состав команды", fmt.Sprintf("chteammemebers|%d", teamID)),
		))
}
