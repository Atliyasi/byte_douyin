package dao

import (
	"gorm.io/gorm"
	"sync"
)

type Message struct {
	gorm.Model
	FromUserId int    `json:"from_user_id"`
	ToUserId   int    `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

func (*Message) TableName() string {
	return "message"
}

type MessageDao struct{}

var messageDao *MessageDao
var messageOnce sync.Once

func NewMessageDao() *MessageDao {
	messageOnce.Do(func() {
		messageDao = &MessageDao{}
	})
	return messageDao
}

// CreateMessage 存储message
func (*MessageDao) CreateMessage(fromUserId int, toUserId int, content string, createTime int64) error {
	message := &Message{
		Model:      gorm.Model{},
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Content:    content,
		CreateTime: createTime,
	}
	if err := GetDB().Create(message).Error; err != nil {
		return err
	}
	return nil
}

// GetMessageById 获取对应的信息列表
func (*MessageDao) GetMessageById(toUserId int, fromUserId int, newTime int64) ([]Message, error) {
	var messageList []Message
	if err := GetDB().Where("(to_user_id=? AND from_user_id=? AND create_time > ?) OR (to_user_id=? AND from_user_id=? AND create_time > ?)", toUserId, fromUserId, newTime, fromUserId, toUserId, newTime).Order("create_time").Find(&messageList).Error; err != nil {
		return nil, err
	}
	return messageList, nil
}
