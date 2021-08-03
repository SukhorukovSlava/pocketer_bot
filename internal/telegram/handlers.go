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
	replyInvalidLink         = "Это невалидная ссылка!"
	replyNotAuthorized       = "Вы не авторизированны! Для авторизации используй комаду /start"
	replyFailAdded           = "Упс, что-то пошло не так, я не смог добавить вашу ссылку!"
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
	msg := tgbotapi.NewMessage(message.Chat.ID, replySuccessfulSavedLink)

	userLink, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = replyInvalidLink
		_, err = b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = replyNotAuthorized
		_, err = b.bot.Send(msg)
		return err
	}

	err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		URL:         userLink.String(),
		AccessToken: accessToken,
	})
	if err != nil {
		msg.Text = replyFailAdded
		_, err = b.bot.Send(msg)
		return err
	}

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
