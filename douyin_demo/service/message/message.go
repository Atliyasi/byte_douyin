package message

import (
	"douyin_demo/dao"
	"sync"
)

type MessageYes struct {
	Id         int    `json:"id"`
	ToUserId   int    `json:"to_user_id"`
	FromUserId int    `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

// SendMessage 实现发送消息逻辑
func SendMessage(fromUserId int, toUserId int, content string, createTime int64) error {
	if err := dao.NewMessageDao().CreateMessage(fromUserId, toUserId, content, createTime); err != nil {
		return err
	}
	return nil
}

// GetMessageList 实现获取正确消息的功能
func GetMessageList(userId int, toUserId int, newTime int64) ([]MessageYes, error) {
	messageList, err := dao.NewMessageDao().GetMessageById(toUserId, userId, newTime)
	if err != nil {
		return nil, err
	}
	var messageYesList []MessageYes
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(messageList))
	for _, message := range messageList {
		go func(message dao.Message) {
			defer wg.Done()
			if message.ToUserId == userId {
				lock.Lock()
				messageYesList = append(messageYesList, MessageYes{
					Id:         int(message.ID),
					ToUserId:   userId,
					FromUserId: toUserId,
					Content:    message.Content,
					CreateTime: message.CreateTime,
				})
				lock.Unlock()
			} else {
				lock.Lock()
				messageYesList = append(messageYesList, MessageYes{
					Id:         int(message.ID),
					ToUserId:   toUserId,
					FromUserId: userId,
					Content:    message.Content,
					CreateTime: message.CreateTime,
				})
				lock.Unlock()
			}

		}(message)
	}
	wg.Wait()
	return messageYesList, nil
}
