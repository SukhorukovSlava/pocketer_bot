package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidUrl      = errors.New("url is invalid")
	errUnauthorized    = errors.New("user is not authorized")
	errUnableToAddLink = errors.New("unable to add link")
)

const (
	invalidLink     = "Это невалидная ссылка!"
	unauthorized    = "Вы не авторизированны! Для авторизации используй комаду /start"
	unableToAddLink = "Упс, что-то пошло не так, я не смог добавить вашу ссылку!"
)

func (b *Bot) handleError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, "")
	switch err {
	case errInvalidUrl:
		msg.Text = invalidLink
	case errUnauthorized:
		msg.Text = unauthorized
	case errUnableToAddLink:
		msg.Text = unableToAddLink
	default:
		msg.Text = "Произошла неизвестная ошибка"
	}
	b.bot.Send(msg)
}
