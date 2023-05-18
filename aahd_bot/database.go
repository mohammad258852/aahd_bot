package aahd_bot

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	Id   int64 `gorm:"primarykey"`
	Name string
}

type AhhdEvent struct {
	MessageId int64
	Date      datatypes.Date `gorm:"primarykey"`
	GroupId   int64
	Group     Group
}

type Status struct {
	gorm.Model
	UserId int64
	User   User
	Read   bool `gorm:"default:false"`
	AhhdId int64
	Ahhd   AhhdEvent
}

type Group struct {
	Id    int64  `gorm:"primarykey"`
	Users []User `gorm:"many2many:user_group;"`
	Name  string
}

var db *gorm.DB

func InitDatabase() error {
	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&AhhdEvent{})
	db.AutoMigrate(&Status{})
	db.AutoMigrate(&Group{})

	return err
}

func GetAhhdEventByDate(group *Group, t time.Time) *AhhdEvent {
	var result AhhdEvent
	err := db.
		Where("date = ? AND group_id = ?", datatypes.Date(t), group.Id).
		First(&result).
		Error
	if err != nil {
		return nil
	}
	return &result
}

func GetAahdEventByMessageId(messageId int64) *AhhdEvent {
	var aahdEvent *AhhdEvent
	err := db.Where("message_id = ?", messageId).First(aahdEvent).Error
	if err != nil {
		return nil
	}
	return aahdEvent
}

func GetAllGroups() []Group {
	var groups []Group

	db.Find(&groups)
	return groups
}

func GetUserStatus(user *User, ahhdEvent *AhhdEvent) *Status {
	var status *Status
	err := db.
		Where("user_id = ? AND ahhd_id = ?", user.Id, ahhdEvent.MessageId).
		First(status).
		Error

	if err != nil {
		return nil
	}
	return status
}

func AddAahdEvent(messageId int64, t time.Time, group *Group) {
	db.Create(&AhhdEvent{messageId, datatypes.Date(t), group.Id, *group})
}

func SaveGroup(group *Group) {
	db.Save(group)
}

func SaveUser(user *User) {
	db.Save(user)
}

func SaveStatus(status *Status) {
	db.Save(status)
}

func GetGroup(id int64) *Group {
	var group *Group
	err := db.First(group, id).Error
	if err != nil {
		return nil
	}
	return group
}

func GetUser(id int64) *User {
	var user *User
	err := db.First(user, id).Error
	if err != nil {
		return nil
	}
	return user
}
