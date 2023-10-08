package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/feof1l/TelegramPocketBot/pkg/repository"
	"github.com/feof1l/TelegramPocketBot/pkg/repository/boltdb"
	"github.com/feof1l/TelegramPocketBot/pkg/server"
	"github.com/feof1l/TelegramPocketBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6390640232:AAGOQzdfgyoDMrZP5HQCVMh5stpKrIfGiM8")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	pocketClient, err := pocket.NewClient("109064-cf3256581e5b77aad9003ab")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/PocetFeof1lBot")
	//telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "https://t.me/PocetFeof1lBot")

	//authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "http://localhost/")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := authorizationServer.Start(); err != nil {

		log.Fatal(err)
	}

}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
