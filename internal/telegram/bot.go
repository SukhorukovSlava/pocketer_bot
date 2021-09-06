package telegram

import (
	"log"

	"github.com/SukhorukovSlava/pocketer_bot/internal/config"
	"github.com/SukhorukovSlava/pocketer_bot/internal/repository"
	"github.com/SukhorukovSlava/pocketer_bot/pkg/pocket"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
	messages        config.Messages
}

func NewBot(
	bot *tgbotapi.BotAPI,
	pocketClient *pocket.Client,
	tr repository.TokenRepository,
	cfg *config.Config,
) *Bot {
	return &Bot{
		bot:             bot,
		pocketClient:    pocketClient,
		tokenRepository: tr,
		redirectUrl:     cfg.AuthServerURL,
		messages:        cfg.Messages,
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	return b.bot.GetUpdatesChan(updateCfg)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		message := update.Message
		if message == nil {
			continue
		}
		if message.IsCommand() {
			if err := b.handleCommand(message); err != nil {
				b.handleError(message.Chat.ID, err)
			}
			continue
		}
		if err := b.handleMessage(message); err != nil {
			b.handleError(message.Chat.ID, err)
		}
	}
}
