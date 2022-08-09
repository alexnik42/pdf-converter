package main

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func startBot() {
	settings := tele.Settings{
		Token:  "5525053757:AAGucYfpFKLLyTWwsBlYtkW67FtXNIYWqvk",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(c tele.Context) error {
		err := c.Send("Hi! I am a pdf converter - upload a file and get it in PDF format")
		if err != nil {
			log.Fatal(err)
		}
		return err
	})

	bot.Handle(tele.OnDocument, convertToPdf)
	bot.Handle(tele.OnPhoto, convertToPdf)

	log.Println("Starting bot")
	bot.Start()
	log.Println("Stopping bot")
}
