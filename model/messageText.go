package model

import (
	"time"
)


//  MType类型有Private(私信)、Public(公共消息)、Global(系统消息)
type MessageText struct {
	ID       uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	SendId   uint64    `json:"sendId" gorm:"column:send_id;not null"`
	Title    string    `json:"title" gorm:"column:title;not null"`
	Text     string    `json:"text" gorm:"column:text;not null"`
	MType    string    `json:"messageType" gorm:"column:message_type;not null"`
	Group    uint64    `json:"groupId" grom:"column:group_id;not null" `
	Role    uint64    `json:"roleId" grom:"column:role_id;not null" `
	PostDate time.Time `gorm:"column:postDate" json:"postDate"`
}

const MessageTextTableName = "tb_message_text"

// TableName :
func (m *MessageText) TableName() string {
	return MessageTextTableName
}

// Create : Create a new MessageText
func (m *MessageText) Create() error {
	m.PostDate = time.Now()
	err := DB.Self.Create(&m).Error
	return err
}


func DeleteMessageText(id uint64) error {

	text := MessageText{}
	text.ID = id
	tx := DB.Self.Begin()
	if err := tx.Delete(&text).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

