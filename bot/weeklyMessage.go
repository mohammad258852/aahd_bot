package bot

import (
	"aahd_bot/db"
	"aahd_bot/util"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageEveryWeek() {
	d := util.GetDurationTilNextWeekDay(time.Friday)

	for {
		time.Sleep(d)
		sendReminderMessageAllGroups()
		d = 7 * 24 * time.Hour
	}
}

func sendReminderMessageAllGroups() {
	for _, group := range db.GetAllGroups() {
		sendReminderMessageToGroup(group)
	}
}

func sendReminderMessageToGroup(group db.Group) {
	t := util.GetCurrentLocalTime()
	text := util.GetReminderText(&group, t, true)

	msg := tgbotapi.NewMessage(group.Id, text)
	msg.ParseMode = "MarkdownV2"
	_, err := bot.Send(msg)
	if err != nil {
		log.Print(err)
	}
}
