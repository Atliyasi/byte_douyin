package comment

import (
	"douyin_demo/dao"
	"strconv"
)

type CommentData struct {
	VideoId     int
	ActionType  string
	CommentText string
	CommentId   string
}

func NewComment(videoId int, actionType string, text string) *CommentData {
	var commentData *CommentData
	commentData = &CommentData{
		VideoId:     videoId,
		ActionType:  actionType,
		CommentText: text,
		CommentId:   text,
	}
	return commentData
}

// AppendComment 实现添加评论逻辑
func (c *CommentData) AppendComment(userId int) (*dao.CommentList, error) {
	comment, err := dao.NewCommentDao().CreateComment(userId, c.VideoId, c.CommentText)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

// DeleteComment 实现删除评论逻辑
func (c *CommentData) DeleteComment() error {
	id, _ := strconv.Atoi(c.CommentId)
	if err := dao.NewCommentDao().DeleteCommentById(id); err != nil {
		return err
	}
	return nil
}

// GetCommentList 实现获取对应视频评论逻辑
func (c *CommentData) GetCommentList() ([]dao.CommentList, error) {
	commentList, err := dao.NewCommentDao().GetCommentListByVideoId(c.VideoId)
	if err != nil {
		return nil, err
	}
	return commentList, nil
}
