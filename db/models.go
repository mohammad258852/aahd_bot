package db

import "gorm.io/datatypes"

type User struct {
	Id   int64 `gorm:"primarykey"`
	Name string
}

type AahdEvent struct {
	MessageId int64 `gorm:"primarykey"`
	Date      datatypes.Date
	GroupId   int64 `gorm:"primarykey"`
	Group     Group
}

type Status struct {
	UserId        int64 `gorm:"primarykey"`
	User          User
	Read          bool  `gorm:"default:false"`
	AhhdMessageId int64 `gorm:"primarykey"`
	AhhdGroupId   int64 `gorm:"primarykey"`
	Ahhd          AahdEvent
}

type Group struct {
	Id    int64  `gorm:"primarykey"`
	Users []User `gorm:"many2many:user_group;"`
	Name  string
}
