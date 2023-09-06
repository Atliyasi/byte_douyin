package controller

import (
	"douyin_demo/dao"
	"douyin_demo/log"
	"douyin_demo/service/comment"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentResponse struct {
	Response
	Comment dao.CommentList `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []dao.CommentList `json:"comment_list"`
}

func HandleComment(c *gin.Context) {
	userIdA, _ := c.Get("uid")
	userId := userIdA.(int)
	videoIdS := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdS)
	actionType := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentId := c.Query("comment_id")
	if actionType == "1" {
		thisComment, err := comment.NewComment(videoId, actionType, commentText).AppendComment(userId)
		if err != nil {
			log.SetLog().Error(err.Error())
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
				Comment: dao.CommentList{},
			})
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "成功",
			},
			Comment: dao.CommentList{
				Id:         thisComment.Id,
				User:       dao.VideoUser{},
				Content:    thisComment.Content,
				CreateDate: thisComment.CreateDate,
			},
		})
	} else if actionType == "2" {
		if err := comment.NewComment(videoId, actionType, commentId).DeleteComment(); err != nil {
			log.SetLog().Error(err.Error())
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
				Comment: dao.CommentList{},
			})
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "删除成功",
			},
			Comment: dao.CommentList{},
		})
	} else {
		log.SetLog().Error("actionType != 2 or 1")
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "失败",
			},
			Comment: dao.CommentList{},
		})
	}
}

func GetCommentList(c *gin.Context) {
	videoIdS := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdS)
	commentList, err := comment.NewComment(videoId, "", "").GetCommentList()
	if err != nil {
		c.JSON(http.StatusBadRequest, CommentListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			CommentList: commentList,
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功",
		},
		CommentList: commentList,
	})
}
