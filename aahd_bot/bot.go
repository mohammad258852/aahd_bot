package aahd_bot

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ptime "github.com/yaa110/go-persian-calendar"
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
	var errorText string
	chatId := update.CallbackQuery.Message.Chat.ID
	group := GetGroup(chatId)
	if group == nil {
		errorText = "گروه در دیتابیس وجود ندارد"
	}

	messageId := update.CallbackQuery.Message.MessageID
	aahdEvent := GetAahdEventByMessageId(int64(messageId))

	if aahdEvent == nil {
		errorText = "پیام در دیتابیس وجود ندارد"
	}

	userId := update.CallbackQuery.From.ID
	user := GetUser(userId)

	if user == nil {
		errorText = "کاربر در دیتابیس وجود ندارد"
	}
	if errorText != "" {
		msg := tgbotapi.NewCallback(update.CallbackQuery.ID, errorText)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}
	data := update.CallbackQuery.Data
	read := data == "1"

	status := GetUserStatus(user, aahdEvent)
	if status == nil {
		status = &Status{User: *user, Ahhd: *aahdEvent, Read: read}
	}
	status.Read = read
	SaveStatus(status)
	updateMessage(group, aahdEvent)

	msg := tgbotapi.NewCallback(update.CallbackQuery.ID, "حله")
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func updateMessage(group *Group, aahdEvent *AhhdEvent) {
	text := getText(group, aahdEvent, time.Time(aahdEvent.Date))
	msg := tgbotapi.NewEditMessageTextAndMarkup(group.Id, int(aahdEvent.MessageId), text, numericKeyboard)
	if _, err := bot.Request(msg); err != nil {
		log.Print(err)
	}
}

func handleMessage(update *tgbotapi.Update) {
	var text string

	switch update.Message.Text {
	case "/in":
		addUser(update)
		text = "خوش اومدی"
	case "/out":
		text = "حیف شد"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func addUser(update *tgbotapi.Update) {
	userId := update.Message.From.ID
	userName := update.Message.From.FirstName + " " + update.Message.From.LastName
	user := &User{Id: userId, Name: userName}
	SaveUser(user)

	chatId := update.Message.Chat.ID
	chatName := update.Message.Chat.Title
	group := &Group{Id: chatId, Name: chatName}
	group.Users = append(group.Users, *user)
	SaveGroup(group)

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
		log.Print(err)
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
		sendMessageAllGroups()
		time.Sleep(d)
		d = 24 * time.Hour
	}
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅", "1"),
		tgbotapi.NewInlineKeyboardButtonData("⛔️", "0"),
	),
)

func messageExist(group *Group, t time.Time) bool {
	ahhd := GetAhhdEventByDate(group, t)
	return ahhd != nil
}

func sendMessageAllGroups() {
	t := time.Now().In(LoadTehranTime())
	for _, group := range GetAllGroups() {
		sendMessageToGroup(group, t)
	}
}

func sendMessageToGroup(group Group, t time.Time) {
	if messageExist(&group, t) {
		return
	}
	text := getText(&group, nil, t)

	msg := tgbotapi.NewMessage(group.Id, text)
	msg.ReplyMarkup = numericKeyboard
	res, err := bot.Send(msg)
	if err != nil {
		log.Print(err)
	}

	AddAahdEvent(int64(res.MessageID), t, &group)
}

func getText(group *Group, aahdEvent *AhhdEvent, t time.Time) string {
	text := group.Name + ":\n"
	p := ptime.New(t)
	text += fmt.Sprintf("🗓 %s/ %d %s %d\n", p.Weekday(), p.Day(), p.Month(), p.Year())

	for _, user := range group.Users {
		text += user.Name + ":" + getStatusString(&user, aahdEvent) + "\n"
	}
	return text
}

func getStatusString(user *User, aahdEvent *AhhdEvent) string {
	if aahdEvent == nil {
		return ""
	}
	status := GetUserStatus(user, aahdEvent)
	if status == nil {
		return ""
	}
	if status.Read {
		return "✅"
	}
	return "⛔️"
}
