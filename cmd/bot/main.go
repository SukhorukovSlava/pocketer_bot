package main

import (
	"log"

	"github.com/SukhorukovSlava/pocketer_bot/internal/config"
	"github.com/SukhorukovSlava/pocketer_bot/internal/repository"
	"github.com/SukhorukovSlava/pocketer_bot/internal/repository/boltdb"
	"github.com/SukhorukovSlava/pocketer_bot/internal/server"
	"github.com/SukhorukovSlava/pocketer_bot/internal/telegram"
	"github.com/SukhorukovSlava/pocketer_bot/pkg/pocket"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := initDb(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(db *bolt.DB) {
		if err = db.Close(); err != nil {
			log.Fatalln(err)
		}
	}(db)

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg)

	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	if err = authorizationServer.Start(); err != nil {
		log.Fatalln(err)
	}
}

func initDb(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		for _, b := range repository.GetBuckets() {
			if _, err = tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
