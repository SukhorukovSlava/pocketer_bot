package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"pocketerClient/pkg/pocket"
)

const (
	cmdStart = "start"

	replyStartTmpl = "Здаствуйте! Чтобы сохранить ссылки в своём Pocket аккаунте, для начала вам необходимо " +
		"придоставить мне доступ. Для этого перейдите по ссылке:\n%s"
	replyAlreadyAuthorized   = "Вы уже авторизованны! Можете присылать ссылку, а я ее сохраню."
	replySuccessfulSavedLink = "Ссылка успешно сохранена!"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.handleStartCmd(message)
	default:
		return b.handleUnknownCmd(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	userLink, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidUrl
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}

	err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		URL:         userLink.String(),
		AccessToken: accessToken,
	})
	if err != nil {
		return errUnableToAddLink
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replySuccessfulSavedLink)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCmd(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorized)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCmd(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введена неизветсная комманда!")
	_, err := b.bot.Send(msg)

	return err
}
