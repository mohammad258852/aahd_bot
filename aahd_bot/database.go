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
	Date   datatypes.Date
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
	res := db.First(&result, "date = ? AND groupId = ?", datatypes.Date(t), group.Id)
	if res.Error != nil {
		return nil
	}
	return &result
}

func GetAllGroups() []Group {
	var groups []Group

	db.Find(&groups)
	return groups
}

func GetUserStatus(user *User, t time.Time, g *Group) *Status {
	var status *Status
	db.First(status, "user = ? AND date = ? AND group = ?", user, datatypes.Date(t), g)

	if status == nil {
		return nil
	}
	return status
}

func AddAahdEvent(messageId int64, t time.Time, group *Group) {
	db.Create(&AhhdEvent{messageId, datatypes.Date(t), group.Id, *group})
}
