package dao

import (
	"douyin_demo/log"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Video struct {
	gorm.Model
	Author        int64  `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

func (*Video) TableName() string {
	return "video"
}

type VideoDao struct{}

var videoDao *VideoDao
var VideoOnce sync.Once

func NewVideDao() *VideoDao {
	VideoOnce.Do(func() {
		videoDao = &VideoDao{}
	})
	return videoDao
}

// GetVideoUser 获取视频发布者信息
func (*VideoDao) GetVideoUser(id int64) *VideoUser {
	var videoUser VideoUser
	if err := GetDB().Where("id=?", id).Find(&videoUser).Error; err != nil {
		log.SetLog().Error("[video_dao] GetVideoUser")
		return nil
	}
	return &videoUser
}

// GetVideoList 获取视频列表
func (this *VideoDao) GetVideoList() ([]Video, error) {
	var video []Video
	if err := GetDB().Find(&video).Error; err != nil {
		log.SetLog().Error("[video_dao] GetVideoList")
		return nil, err
	}
	return video, nil
}

// GetVideoListById 通过用户Id获取对应发布的视频信息
func (this *VideoDao) GetVideoListById(id int) ([]Video, error) {
	var video []Video
	if err := GetDB().Find(&video, "author=?", id).Error; err != nil {
		log.SetLog().Error("[video_dao] GetVideoListById")
		return nil, err
	}
	return video, nil
}

// GetVideoById 通过Id获取视频信息
func (*VideoDao) GetVideoById(id int64) *Video {
	var videoList Video
	err := GetDB().Where("id=?", id).First(&videoList).Error
	if err != nil {
		log.SetLog().Error("[video_dao] GetVideoById")
		return nil
	}
	return &videoList
}

// SetVideo 创建视频数据
func (*VideoDao) SetVideo(video *Video) error {
	if err := GetDB().Create(video).Error; err != nil {
		log.SetLog().Error("[video_dao] SetVideo")
		return err
	}
	return nil
}

// QueryVideoListByTime 通过时间一共len条获取最新的视频数据
func (*VideoDao) QueryVideoListByTime(lastTime time.Time, len int) (*[]Video, error) {
	var videos []Video
	if err := GetDB().Where("updated_time < ?", lastTime).Order("updated_time desc").Limit(len).Find(&videos).Error; err != nil {
		log.SetLog().Error("[video_dao] QueryVideoListByTime")
		return nil, err
	}
	return &videos, nil
}

// UpdateVideoUserByWork 更新用户发布视频个数+1
func (*VideoDao) UpdateVideoUserByWork(id int64) error {
	if err := GetDB().Model(&VideoUser{}).Where("id=?", id).Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
		log.SetLog().Error("[video_dao] UpdateVideoUserByWork")
		return err
	}
	return nil
}
