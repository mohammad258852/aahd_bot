package bot

import (
	"aahd_bot/db"
	"aahd_bot/util"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func CreateBot() error {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		return err
	}

	bot.Debug = true

	return err
}

func RunBot() {
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := getUpdateChannel()

	for update := range updates {
		handleUpdate(&update)
	}
}

func updateMessage(group *db.Group, aahdEvent *db.AahdEvent) {
	text := util.GetText(group, aahdEvent, time.Time(aahdEvent.Date), true)
	msg := tgbotapi.NewEditMessageTextAndMarkup(group.Id, int(aahdEvent.MessageId), text, util.NumericKeyboard)
	msg.ParseMode = "MarkdownV2"
	if _, err := bot.Request(msg); err != nil {
		log.Print(err)
	}
}

func getUpdateChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"chat_member", "message", "inline_query", "callback_query"}

	updates := bot.GetUpdatesChan(u)
	return updates
}
