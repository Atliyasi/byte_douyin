package user

import (
	"douyin_demo/dao"
	"douyin_demo/util"
	"errors"
	"regexp"
)

type UserLoginFlow struct {
	username string
	password string
}

// NewUserLoginFlow 入口函数
func NewUserLoginFlow(username string, pwd string) *UserLoginFlow {
	return &UserLoginFlow{
		username: username,
		password: pwd,
	}
}

// Login 封装Login步骤
func Login(username string, pwd string) (int, error) {
	return NewUserLoginFlow(username, pwd).Do()
}

// Do 使用UserLoginFlow的Login函数进行登陆
func (u *UserLoginFlow) Do() (int, error) {
	if err := u.CheckParam(); err != nil {
		return 0, err
	}
	id, err := u.Login()
	if err != nil {
		return 0, err
	}
	return id, nil

}

// CheckParam 用户名和密码规范逻辑
func (u *UserLoginFlow) CheckParam() error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	if reg.MatchString(u.username) {
		return nil
	}
	return errors.New("请输入正确的邮箱")
}

// Login 用户登陆逻辑
func (u *UserLoginFlow) Login() (int, error) {
	//判断用户是否存在
	userQuery, err := dao.NewUserDao().FindUserByName(u.username)
	if err != nil || userQuery.ID == 0 {
		return 0, errors.New("用户不存在")
	}
	//判断密码是否正确
	if !util.PwdVerify(u.password, userQuery.Password) {
		return 0, errors.New("密码错误")
	}
	//登录状态
	return int(userQuery.ID), nil
}
