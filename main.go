package main

import (
	"aahd_bot/bot"
	"aahd_bot/db"
	"log"
)

func main() {
	err := db.InitDatabase()
	if err != nil {
		log.Print(err)
	}
	err = bot.CreateBot()
	if err != nil {
		log.Printf("Error in Creating Bot: %s", err)
		return
	}

	go bot.SendMessageEveryDay()
	bot.RunBot()
}
