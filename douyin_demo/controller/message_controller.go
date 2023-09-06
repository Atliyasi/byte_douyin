package controller

import (
	"douyin_demo/log"
	"douyin_demo/service/message"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type MessageResponse struct {
	Response
	MessageList []message.MessageYes `json:"message_list"`
}

func SendMessage(c *gin.Context) {
	now := time.Now().Unix()
	fromUserIdA, _ := c.Get("uid")
	toUserIdS := c.Query("to_user_id")
	fromUserId := fromUserIdA.(int)
	toUserId, _ := strconv.Atoi(toUserIdS)
	actionType := c.Query("action_type")
	content := c.Query("content")
	if actionType == "1" {
		if err := message.SendMessage(fromUserId, toUserId, content, now); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {
		log.SetLog().Error("actionType != \"1\"")
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "失败",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "发送成功",
	})
}

func GetMessage(c *gin.Context) {
	userIdA, _ := c.Get("uid")
	toUserIdS := c.Query("to_user_id")
	times := c.Query("pre_msg_time")
	userId := userIdA.(int)
	toUserId, _ := strconv.Atoi(toUserIdS)
	msgTime, _ := strconv.Atoi(times)
	newTime := int64(msgTime)
	messageYesList, err := message.GetMessageList(userId, toUserId, newTime)
	if err != nil {
		log.SetLog().Error(err.Error())
		c.JSON(http.StatusBadRequest, MessageResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			MessageList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功",
		},
		MessageList: messageYesList,
	})
}
