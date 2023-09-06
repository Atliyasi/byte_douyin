package relation

import "douyin_demo/dao"

// FriendList 返回朋友列表
func FriendList(userId int) ([]dao.VideoUser, error) {
	VideoUserList, err := dao.NewRelationDao().FriendList(userId)
	if err != nil {
		return nil, err
	}
	return VideoUserList, nil
}
