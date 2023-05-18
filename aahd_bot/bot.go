package aahd_bot

import (
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

func handleUpdate(update *tgbotapi.Update) {
	if update.Message != nil {
		handleMessage(update)
	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update)
	}
}

func handleCallbackQuery(update *tgbotapi.Update) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func handleMessage(update *tgbotapi.Update) {
	var text string

	switch update.Message.Text {
	case "/in":
		text = "خوش اومدی"
	case "/out":
		text = "حیف شد"
	default:
		text = "متوجه نشدم"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func getUpdateChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"chat_member", "message", "inline_query", "callback_query"}

	updates := bot.GetUpdatesChan(u)
	return updates
}

var tehranTime *time.Location

func LoadTehranTime() *time.Location {
	if tehranTime != nil {
		return tehranTime
	}
	var err error
	tehranTime, err = time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Panic(err)
	}
	return tehranTime
}

func SendMessageEveryDay() {
	t := time.Now().In(LoadTehranTime())
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 1, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}

	for {
		if !messageExist() {
			sendMessage()
		}
		time.Sleep(d)
		d = 24 * time.Hour
	}
}

func messageExist() bool {
	ahhd := GetAhhdEventByDate(time.Now().In(LoadTehranTime()))
	return ahhd != nil
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("خوندم", "1"),
		tgbotapi.NewInlineKeyboardButtonData("صدقه", "0"),
	),
)

func sendMessage() {
	t := time.Now()
	for _, group := range GetAllGroups() {
		text := group.Name + "\n"
		for _, user := range group.Users {
			text += user.Name + ":" + getStatusString(&user, t, &group) + "\n"
		}

		msg := tgbotapi.NewMessage(group.Id, text)
		msg.ReplyMarkup = numericKeyboard
		res, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

		AddAahdEvent(int64(res.MessageID), t, &group)
	}
}

func getStatusString(user *User, t time.Time, g *Group) string {
	status := GetUserStatus(user, t, g)
	if status == nil {
		return ""
	}
	if status.Read {
		return "✅"
	}
	return "⛔️"
}
