package controller

import (
	"douyin_demo/log"
	"douyin_demo/middleware"
	"douyin_demo/service/publish"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []publish.VideoList `json:"video_list,omitempty"`
	NextTime  int64               `json:"next_time,omitempty"`
}

func PublishAction(c *gin.Context) {
	token := c.PostForm("token")
	_, claims, err := middleware.ParseToken(token)
	id := claims.UserId
	data, err := c.FormFile("data")
	if err != nil {
		log.SetLog().Error("[PublishAction]" + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", id, filename)
	saveFile := filepath.Join("../public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		log.SetLog().Error("[PublishAction]" + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  finalName + "上传失败",
		})
		return
	}
	title := c.PostForm("title")
	if err := publish.SavePublishVideo(saveFile, finalName, title, int64(id)); err != nil {
		log.SetLog().Error("[PublishAction] " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  finalName + "上传失败",
		})
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + "上传成功",
	})
	log.SetLog().Info("[PublishAction] " + strconv.FormatInt(int64(id), 10) + " " + finalName)
}

// Feed 获取静态视频资源信息
func Feed(c *gin.Context) {
	token := c.Query("token")
	_, claims, _ := middleware.ParseToken(token)
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "请求成功",
		},
		VideoList: publish.AssembleVideoList(claims.UserId),
		NextTime:  time.Now().Unix(),
	})
}

func PublishList(c *gin.Context) {
	userId := c.Query("user_id")
	id, _ := strconv.Atoi(userId)
	//fmt.Println("id: ", id)
	if id == 0 {
		token := c.Query("token")
		_, claims, err := middleware.ParseToken(token)
		if err != nil {
			return
		}
		id = claims.UserId
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "请求成功",
		},
		VideoList: publish.AssembleVideoListById(id),
		NextTime:  time.Now().Unix(),
	})
}
