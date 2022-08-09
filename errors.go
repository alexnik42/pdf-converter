package main

import (
	"log"

	tele "gopkg.in/telebot.v3"
)

func logErrorEvent(err error, c tele.Context) error {
	log.Printf("[%s]: error - %s\n", c.Message().Sender.Username, err)
	return err
}

func logInfoEvent(event string, c tele.Context) {
	log.Printf("[%s]: %s\n", c.Message().Sender.Username, event)
}
