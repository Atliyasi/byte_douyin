package controller

import (
	"douyin_demo/log"
	"douyin_demo/service/favorite"
	"douyin_demo/service/publish"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteResponse struct {
	Response
	VideoList []publish.VideoList `json:"video_list,omitempty"`
}

func Upvote(c *gin.Context) {
	videoIdS := c.Query("video_id")
	actionTypeS := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdS)
	actionType, _ := strconv.Atoi(actionTypeS)
	userIdS, _ := c.Get("uid")
	userId, _ := userIdS.(int)
	err := favorite.NewFavoriteFlow(videoId, actionType).Upvote(userId)
	if err != nil {
		log.SetLog().Error(err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "成功",
	})
}

func FavoriteList(c *gin.Context) {
	userIdS := c.Query("user_id")
	userId, _ := strconv.Atoi(userIdS)
	favoriteList, err := favorite.FavoriteList(userId)
	if err != nil {
		log.SetLog().Error(err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, FavoriteResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功",
		},
		VideoList: favoriteList,
	})
}
