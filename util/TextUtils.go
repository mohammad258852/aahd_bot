package util

import (
	"aahd_bot/db"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ptime "github.com/yaa110/go-persian-calendar"
	"regexp"
	"strconv"
	"time"
)

var NumericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅", "1"),
		tgbotapi.NewInlineKeyboardButtonData("⛔️", "0"),
	),
)

func GetText(group *db.Group, aahdEvent *db.AhhdEvent, t time.Time, markdown bool) string {
	text := group.Name + ":\n"
	p := ptime.New(t)
	text += fmt.Sprintf("🗓 %s/ %d %s %d\n", p.Weekday(), p.Day(), p.Month(), p.Year())

	for _, user := range group.Users {
		statusStr := getStatusString(&user, aahdEvent)
		if markdown {

			text += fmt.Sprintf("[%s](tg://user?id=%d):%s\n", escapedMarkdownText(user.Name), user.Id, escapedMarkdownText(statusStr))
		} else {
			text += user.Name + ":" + statusStr + "\n"
		}
	}
	return text
}

func GetReminderText(group *db.Group, t time.Time, markdown bool) string {
	text := group.Name + ":\n"
	text += fmt.Sprintf("یادآوری هفت روز گذشته\n")

	for _, user := range group.Users {
		unreadCount := getUnreadCount(&user, group, t)
		if markdown {
			text += fmt.Sprintf("[%s](tg://user?id=%d):%d روز\n", escapedMarkdownText(user.Name), user.Id, unreadCount)
		} else {
			text += user.Name + ":" + strconv.Itoa(unreadCount) + "\n"
		}
	}
	return text
}

func getUnreadCount(user *db.User, group *db.Group, t time.Time) int {
	read := db.GetWeeklyRead(user, group, t)
	if read == nil {
		return 7
	}
	return 7 - *read
}

var markdownEscapeRegex = regexp.MustCompile(`([.#*_{}\[\]])`)

func escapedMarkdownText(text string) string {
	return markdownEscapeRegex.ReplaceAllString(text, `\$1`)
}

func getStatusString(user *db.User, aahdEvent *db.AhhdEvent) string {
	if aahdEvent == nil {
		return ""
	}
	status := db.GetUserStatus(user, aahdEvent)
	if status == nil {
		return ""
	}
	if status.Read {
		return "✅"
	}
	return "⛔️"
}
