package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func startBot() {
	settings := tele.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(c tele.Context) error {
		logInfoEvent("Activating bot", c)
		err := c.Send("Hi! I am a pdf converter - upload files and get them in PDF format")
		if err != nil {
			log.Println(err)
		}
		return err
	})
	bot.Handle(tele.OnDocument, convertToPDF)
	bot.Handle(tele.OnPhoto, func(c tele.Context) error {
		err := c.Send("Please send me the picture as a 'File', not as a 'Photo'")
		if err != nil {
			log.Println(err)
		}
		return err
	})
	bot.Handle(tele.OnVideo, func(c tele.Context) error {
		err := c.Send("Could't convert 'Video' to PDF format")
		if err != nil {
			log.Println(err)
		}
		return err
	})

	log.Println("Starting bot")
	bot.Start()
	log.Println("Stopping bot")
}
