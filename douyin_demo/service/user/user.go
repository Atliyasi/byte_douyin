package user

import (
	"douyin_demo/dao"
)

func QueryUserById(id int) (*dao.VideoUser, error) {
	user, err := dao.NewUserDao().FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
