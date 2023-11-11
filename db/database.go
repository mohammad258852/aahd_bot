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
	err = db.AutoMigrate(&AahdEvent{})
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

func GetAhhdEventByDate(group *Group, t time.Time) *AahdEvent {
	var result AahdEvent
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

func GetAahdEventByMessageId(messageId int64) *AahdEvent {
	var aahdEvent AahdEvent
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

func GetUserStatus(user *User, ahhdEvent *AahdEvent) *Status {
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

func AddAahdEvent(messageId int64, t time.Time, group *Group) {
	db.Create(&AahdEvent{messageId, datatypes.Date(t), group.Id, *group})
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
