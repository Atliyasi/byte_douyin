package dao

import (
	"douyin_demo/util"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type User struct {
	gorm.Model
	Username  string `gorm:"column:username" gorm:"uniqueIndex"`
	Password  string `gorm:"column:password"`
	IsDeleted bool   `gorm:"column:is_deleted"`
}

type VideoUser struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int    `json:"total_favorited,omitempty"`
	WorkCount       int    `json:"work_count,omitempty"`
	FavoriteCount   int    `json:"favorite_count,omitempty"`
}

func (*User) TableName() string {
	return "user"
}

func (*VideoUser) TableName() string {
	return "video_user"
}

type UserDao struct{}

var userDao *UserDao
var userOnce sync.Once

// NewUserDao 创建UserDao实例
func NewUserDao() *UserDao {
	userOnce.Do(func() {
		userDao = &UserDao{}
	})
	return userDao
}

// FindUserByName 通过user账号查询user信息
func (*UserDao) FindUserByName(username string) (*User, error) {
	var user User
	if err := GetDB().Where("username=?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// RegisterUser 创建用户信息
func (*UserDao) RegisterUser(user *User) (int, error) {
	db := GetDB()
	tx := db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return -1, err
	}
	number := util.RandInt()
	videoUser := VideoUser{
		Id:              int64(user.ID),
		Name:            user.Username,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          fmt.Sprintf("http://%s:%s/%d.png", HOST, PORT, number),
		BackgroundImage: fmt.Sprintf("http://%s:%s/%d.png", HOST, PORT, number),
		Signature:       "内测用户",
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
	if err := tx.Create(videoUser).Error; err != nil {
		tx.Rollback()
		return -1, err
	}
	tx.Commit()
	return int(user.ID), nil
}

// FindUserById 通过userId查询对应user的详细信息
func (*UserDao) FindUserById(id int) (*VideoUser, error) {
	var user VideoUser
	if err := GetDB().Where("id=?", id).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
