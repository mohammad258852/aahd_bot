package bot

import (
	"aahd_bot/db"
	"aahd_bot/util"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageEveryDay() {
	d := util.GetDurationTilNextDay()

	for {
		sendMessageAllGroups()
		time.Sleep(d)
		d = 24 * time.Hour
	}
}

func sendMessageAllGroups() {
	t := util.GetCurrentLocalTime()
	for _, group := range db.GetAllGroups() {
		sendMessageToGroup(group, t)
	}
}

func sendMessageToGroup(group db.Group, t time.Time) {
	if messageExist(&group, t) {
		return
	}
	text := util.GetText(&group, nil, t, true)

	msg := tgbotapi.NewMessage(group.Id, text)
	msg.ReplyMarkup = util.NumericKeyboard
	msg.ParseMode = "MarkdownV2"
	res, err := bot.Send(msg)
	if err != nil {
		log.Print(err)
	}

	db.AddAahdEvent(int64(res.MessageID), t, &group)
}

func messageExist(group *db.Group, t time.Time) bool {
	ahhd := db.GetAhhdEventByDate(group, t)
	return ahhd != nil
}
