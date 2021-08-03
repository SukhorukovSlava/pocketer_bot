package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"pocketerClient/internal/repository"
	"pocketerClient/pkg/pocket"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Start, authLink))

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.tokenRepository.Get(chatId, repository.AccessTokens)
}

func (b *Bot) generateAuthorizationLink(chatId int64) (pocket.AuthorizationUrl, error) {
	redirectUrl := fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatId)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectUrl)
	if err != nil {
		return "", err
	}

	err = b.tokenRepository.Put(chatId, requestToken, repository.RequestTokens)
	if err != nil {
		return "", err
	}

	return b.pocketClient.MakeAuthorizationUrl(requestToken, redirectUrl)
}
