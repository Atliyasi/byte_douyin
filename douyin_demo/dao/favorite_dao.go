package dao

import (
	"errors"
	"gorm.io/gorm"
	"sync"
)

type FavoriteList struct {
	gorm.Model
	UserId     int  `gorm:"column:user_id"`
	VideoId    int  `gorm:"column:video_id"`
	IsFavorite bool `gorm:"column:is_favorite"`
}

func (*FavoriteList) TableName() string {
	return "favorite_list"
}

type FavoriteDao struct{}

var favoriteDao *FavoriteDao
var FavoriteOnce sync.Once

func NewFavoriteDao() *FavoriteDao {
	FavoriteOnce.Do(func() {
		favoriteDao = &FavoriteDao{}
	})
	return favoriteDao
}

// Like 点赞 +1
func (f *FavoriteDao) Like(videoId int, userId int) error {
	db := GetDB()
	tx := db.Begin()
	// 通过视频id查找到对应视频信息，并对喜欢数+1
	var video Video
	if err := tx.First(&video, "id=?", videoId).Error; err != nil {
		tx.Rollback()
		return err
	}
	video.FavoriteCount++
	if err := tx.Save(&video).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 通过视频发布者id获取对应用户信息并且该用户获赞+1
	var videoUser VideoUser
	if err := tx.Where("id=?", video.Author).First(&videoUser).Error; err != nil || videoUser.Id == int64(userId) {
		tx.Rollback() // 回滚事务
		return err
	}
	videoUser.TotalFavorited++
	if err := tx.Save(&videoUser).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 通过用户id获取对应用户信息并且该用户喜欢+1
	var user VideoUser
	if err := tx.Where("id=?", userId).First(&user).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	user.FavoriteCount++
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 记录喜欢信息
	var favoriteList FavoriteList
	if err := tx.Where("user_id=?", userId).Where("video_id=?", videoId).First(&favoriteList).Error; err != nil {
		if err := f.CreateFavoriteList(userId, videoId, true); err != nil {
			return err
		}
	} else {
		if err := f.UpdateFavoriteList(userId, videoId, true); err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// Unlike 取消点赞 -1
func (f *FavoriteDao) Unlike(videoId int, userId int) error {
	db := GetDB()
	// 开始事务
	tx := db.Begin()
	var video Video
	if err := tx.First(&video, "id=?", videoId).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	video.FavoriteCount--
	if err := tx.Save(&video).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	var videoUser VideoUser
	if err := tx.Where("id=?", video.Author).First(&videoUser).Error; err != nil || videoUser.Id == int64(userId) {
		tx.Rollback() // 回滚事务
		return err
	}
	videoUser.TotalFavorited--
	if err := tx.Save(&videoUser).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 通过用户id获取对应用户信息并且该用户喜欢-1
	var user VideoUser
	if err := tx.Where("id=?", userId).First(&user).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("有误")
	}
	user.FavoriteCount--
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	if err := f.UpdateFavoriteList(userId, videoId, false); err != nil {
		return err
	}
	return nil
}

func (*FavoriteDao) CreateFavoriteList(userId int, videoId int, isFavorite bool) error {
	favorite := &FavoriteList{
		Model:      gorm.Model{},
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: isFavorite,
	}
	if err := GetDB().Create(favorite).Error; err != nil {
		return err
	}
	return nil
}

func (*FavoriteDao) UpdateFavoriteList(userId int, videoId int, isFavorite bool) error {
	err := GetDB().Model(&FavoriteList{}).Where("user_id=?", userId).Where("video_id=?", videoId).Update("IsFavorite", isFavorite).Error
	if err != nil {
		return err
	}
	return nil
}

func (*FavoriteDao) FindFavoriteByUserIdByVideoId(userId int, videoId int) bool {
	var favoriteList FavoriteList
	if err := GetDB().Where("user_id=?", userId).Where("video_id=?", videoId).Find(&favoriteList).Error; err != nil {
		return false
	}
	return favoriteList.IsFavorite
}

func (*FavoriteDao) FindFavoriteList(userId int) ([]FavoriteList, error) {
	var favoriteList []FavoriteList
	if err := GetDB().Where("user_id=?", userId).Where("is_favorite=?", true).Find(&favoriteList).Error; err != nil {
		return nil, err
	}
	return favoriteList, nil
}

func (f *FavoriteDao) FindFavoritesByUserId(id int) ([]FavoriteList, error) {
	var favoriteList []FavoriteList
	if err := GetDB().Where("user_id=?", id).Where("is_favorite=?", 1).Find(&favoriteList).Error; err != nil {
		return nil, err
	}
	return favoriteList, nil
}
