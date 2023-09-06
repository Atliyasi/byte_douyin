package dao

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Comment struct {
	gorm.Model
	VideoId    int    `json:"video_id"`
	UserId     int    `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

func (*Comment) TableName() string {
	return "comment"
}

type CommentList struct {
	Id         int64     `json:"id,omitempty"`
	User       VideoUser `json:"user"`
	Content    string    `json:"content,omitempty"`
	CreateDate string    `json:"create_date,omitempty"`
}

type CommentDao struct{}

var commentDao *CommentDao
var CommentOnce sync.Once

func NewCommentDao() *CommentDao {
	CommentOnce.Do(func() {
		commentDao = &CommentDao{}
	})
	return commentDao
}

// CreateComment 创建一条评论
func (*CommentDao) CreateComment(userId int, videoId int, content string) (*CommentList, error) {
	_, M, d := time.Now().Date()
	timeString := fmt.Sprintf("%d-%d", M, d)
	comment := &Comment{
		Model:      gorm.Model{},
		VideoId:    videoId,
		UserId:     userId,
		Content:    content,
		CreateDate: timeString,
	}
	db := GetDB()
	tx := db.Begin()
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var user VideoUser
	if err := tx.Where("id=?", comment.UserId).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var video Video
	if err := tx.Where("id=?", videoId).Find(&video).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	video.CommentCount++
	if err := tx.Save(video).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &CommentList{
		Id:         int64(comment.ID),
		User:       user,
		Content:    content,
		CreateDate: timeString,
	}, nil
}

// DeleteCommentById 根据评论id删除一条评论
func (*CommentDao) DeleteCommentById(commentId int) error {
	var comment Comment
	db := GetDB()
	tx := db.Begin()
	if err := tx.Where("id=?", commentId).Delete(&comment).Error; err != nil {
		return err
	}
	var video Video
	if err := tx.Where("id=?", comment.VideoId).Find(&video).Error; err != nil {
		tx.Rollback()
		return err
	}
	video.CommentCount--
	if err := tx.Save(video).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// GetCommentListByVideoId 通过视频id获取对应视频的评论信息
func (*CommentDao) GetCommentListByVideoId(videoId int) ([]CommentList, error) {
	var comments []Comment
	var commentList []CommentList
	if err := GetDB().Where("video_id=?", videoId).Find(&comments).Error; err != nil {
		return nil, err
	}
	for _, comment := range comments {
		var user VideoUser
		if err := GetDB().Where("id=?", comment.UserId).First(&user).Error; err != nil {
			continue
		}
		commentList = append(commentList, CommentList{
			Id:         int64(comment.ID),
			User:       user,
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		})
	}
	return commentList, nil
}
