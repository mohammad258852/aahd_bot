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
		return
	}
	err = bot.CreateBot()
	if err != nil {
		log.Printf("Error in Creating Bot: %s", err)
		return
	}

	go bot.SendMessageEveryDay()
	go bot.SendMessageEveryWeek()
	bot.RunBot()
}
