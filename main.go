package main

import (
	"aahd_bot/aahd_bot"
	"log"
)

func main() {
	err := aahd_bot.InitDatabase()
	if err != nil {
		log.Print(err)
	}
	err = aahd_bot.CreateBot()
	if err != nil {
		log.Print(err)
	}

	go aahd_bot.SendMessageEveryDay()
	aahd_bot.RunBot()
}
