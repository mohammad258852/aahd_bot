package db

import (
	"log"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func InitDatabase() error {
	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&AhhdEvent{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Status{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Group{})
	if err != nil {
		return err
	}

	return err
}

func GetAhhdEventByDate(group *Group, t time.Time) *AhhdEvent {
	var result AhhdEvent
	err := db.
		Where("date = ? AND group_id = ?", datatypes.Date(t), group.Id).
		First(&result).
		Error
	if err != nil {
		log.Print(err)
		return nil
	}
	return &result
}

func GetAahdEventByMessageId(messageId int64) *AhhdEvent {
	var aahdEvent AhhdEvent
	err := db.Where("message_id = ?", messageId).First(&aahdEvent).Error
	if err != nil {
		log.Print(err)
		return nil
	}
	return &aahdEvent
}

func GetAllGroups() []Group {
	var groups []Group

	db.Preload(clause.Associations).Find(&groups)
	return groups
}

func GetUserStatus(user *User, ahhdEvent *AhhdEvent) *Status {
	var status Status
	err := db.
		Where("user_id = ? AND ahhd_message_id = ? AND ahhd_group_id", user.Id, ahhdEvent.MessageId, ahhdEvent.GroupId).
		First(&status).
		Error

	if err != nil {
		log.Print(err)
		return nil
	}
	return &status
}

func GetWeeklyRead(user *User, group *Group, t time.Time) *int {
	var result int
	err := db.Table("statuses as s").
		Select("count(*)").
		Joins("join ahhd_events as a on s.ahhd_group_id = a.group_id and s.ahhd_message_id = a.message_id").
		Where("s.user_id = ? AND s.ahhd_group_id = ? AND a.date < datetime(?, '-1 day') AND a.date > datetime(?, '-8 day')", user.Id, group.Id, t, t).
		Scan(&result).
		Error

	if err != nil {
		log.Print(err)
		return nil
	}
	return &result
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
	var group Group
	err := db.Preload(clause.Associations).First(&group, id).Error
	if err != nil {
		log.Print(err)
		return nil
	}
	return &group
}

func GetUser(id int64) *User {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		log.Print(err)
		return nil
	}
	return &user
}

func DeleteUserFromGroup(userId int64, groupId int64) {
	group := GetGroup(groupId)
	user := GetUser(userId)

	if group == nil || user == nil {
		return
	}
	err := db.Model(group).Association("Users").Delete(user)
	if err != nil {
		log.Printf("error in deleting user %s from group %s. error: %s", user.Name, group.Name, err)
		return
	}
}
