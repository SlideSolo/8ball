package main

import (
	"./database"
	"github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"time"
)


func main()  {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))

	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		err = database.CollectData(int(update.Message.Chat.ID), update.Message.Chat.UserName, update.Message.Text)
		if err != nil {
			log.Fatal(err)
		}
		switch update.Message.Text {
		case "/start":

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a Magic 8Ball. Write your question below")
			bot.Send(msg)
		default:

			rand.Seed(time.Now().UnixNano())
			phraseIndex := rand.Intn(19)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, database.GetAnswer(phraseIndex) )
			bot.Send(msg)
		}

	}
}




