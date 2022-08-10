package main

import (
	"log"

	tele "gopkg.in/telebot.v3"
)

func logErrorEvent(err error, c tele.Context) {
	log.Printf("[%s]: file - [%s], error - %s\n", c.Message().Sender.Username, c.Message().Document.FileName, err)
}

func logInfoEvent(event string, c tele.Context) {
	log.Printf("[%s]: file - [%s], status - %s\n", c.Message().Sender.Username, c.Message().Document.FileName, event)
}
