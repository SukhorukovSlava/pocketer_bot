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

func (b *Bot) handleError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, "")
	switch err {
	case errInvalidUrl:
		msg.Text = b.messages.Errors.InvalidLink
	case errUnauthorized:
		msg.Text = b.messages.Errors.Unauthorized
	case errUnableToAddLink:
		msg.Text = b.messages.Errors.UnableToAddLink
	default:
		msg.Text = b.messages.Errors.Default
	}
	b.bot.Send(msg)
}
