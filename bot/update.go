package bot

import (
	"aahd_bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
)

func handleUpdate(update *tgbotapi.Update) {
	if update.Message != nil {
		handleMessage(update)
	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update)
	}
}

func handleMessage(update *tgbotapi.Update) {
	var text *string

	switch update.Message.Text {
	case "/in":
		text = addUser(update)
	case "/out":
		text = deleteUser(update)
	}

	if r.MatchString(update.Message.Text) {
		text = rename(update)
	}

	if text != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, *text)
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
	}

	if update.Message.LeftChatMember != nil {
		handleLefChat(update)
	}
}

func handleCallbackQuery(update *tgbotapi.Update) {
	var errorText string
	chatId := update.CallbackQuery.Message.Chat.ID
	group := db.GetGroup(chatId)
	if group == nil {
		errorText = "گروه در دیتابیس وجود ندارد"
	}

	messageId := update.CallbackQuery.Message.MessageID
	aahdEvent := db.GetAahdEventByMessageId(int64(messageId))

	if aahdEvent == nil {
		errorText = "پیام در دیتابیس وجود ندارد"
	}

	userId := update.CallbackQuery.From.ID
	user := db.GetUser(userId)

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

	status := db.GetUserStatus(user, aahdEvent)
	if status == nil {
		status = &db.Status{User: *user, Ahhd: *aahdEvent, Read: read}
	}
	status.Read = read
	db.SaveStatus(status)
	updateMessage(group, aahdEvent)

	msg := tgbotapi.NewCallback(update.CallbackQuery.ID, "حله")
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

var r, _ = regexp.Compile(`/name\w*\s+(.*)`)

func deleteUser(update *tgbotapi.Update) *string {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	db.DeleteUserFromGroup(userId, chatId)
	text := "حیف شد"
	return &text
}

func handleLefChat(update *tgbotapi.Update) {
	userId := update.Message.LeftChatMember.ID
	chatId := update.Message.Chat.ID
	db.DeleteUserFromGroup(userId, chatId)
}

func rename(update *tgbotapi.Update) *string {
	userId := update.Message.From.ID
	userName := r.FindStringSubmatch(update.Message.Text)[1]
	user := &db.User{Id: userId, Name: userName}
	db.SaveUser(user)
	text := "حله"
	return &text
}

func addUser(update *tgbotapi.Update) *string {
	userId := update.Message.From.ID
	userName := update.Message.From.FirstName + " " + update.Message.From.LastName
	user := &db.User{Id: userId, Name: userName}
	db.SaveUser(user)

	chatId := update.Message.Chat.ID
	group := db.GetGroup(chatId)
	if group == nil {
		chatName := update.Message.Chat.Title
		group = &db.Group{Id: chatId, Name: chatName}
	}
	for _, u := range group.Users {
		if u.Id == user.Id {
			return nil
		}
	}
	group.Users = append(group.Users, *user)
	db.SaveGroup(group)
	text := "خوش اومدی"
	return &text
}
